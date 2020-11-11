package adapters

import (
	"context"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/app/query"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TrainingModel struct {
	UUID     string `firestore:"Uuid"`
	UserUUID string `firestore:"UserUuid"`
	User     string `firestore:"User"`

	Time  time.Time `firestore:"Time"`
	Notes string    `firestore:"Notes"`

	ProposedTime   *time.Time `firestore:"ProposedTime"`
	MoveProposedBy *string    `firestore:"MoveProposedBy"`

	Canceled bool `firestore:"Canceled"`
}

type TrainingsFirestoreRepository struct {
	firestoreClient *firestore.Client
}

func NewTrainingsFirestoreRepository(
	firestoreClient *firestore.Client,
) TrainingsFirestoreRepository {
	return TrainingsFirestoreRepository{
		firestoreClient: firestoreClient,
	}
}

func (r TrainingsFirestoreRepository) trainingsCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("trainings")
}

func (r TrainingsFirestoreRepository) AddTraining(ctx context.Context, tr *training.Training) error {
	collection := r.trainingsCollection()

	trainingModel := r.marshalTraining(tr)

	return r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return tx.Create(collection.Doc(trainingModel.UUID), trainingModel)
	})
}

func (r TrainingsFirestoreRepository) GetTraining(
	ctx context.Context,
	trainingUUID string,
	user training.User,
) (*training.Training, error) {
	firestoreTraining, err := r.trainingsCollection().Doc(trainingUUID).Get(ctx)

	if status.Code(err) == codes.NotFound {
		return nil, training.NotFoundError{trainingUUID}
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to get actual docs")
	}

	tr, err := r.unmarshalTraining(firestoreTraining)
	if err != nil {
		return nil, err
	}

	if err := training.CanUserSeeTraining(user, *tr); err != nil {
		return nil, err
	}

	return tr, nil
}

func (r TrainingsFirestoreRepository) UpdateTraining(
	ctx context.Context,
	trainingUUID string,
	user training.User,
	updateFn func(ctx context.Context, tr *training.Training) (*training.Training, error),
) error {
	trainingsCollection := r.trainingsCollection()

	return r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		documentRef := trainingsCollection.Doc(trainingUUID)

		firestoreTraining, err := tx.Get(documentRef)
		if err != nil {
			return errors.Wrap(err, "unable to get actual docs")
		}

		tr, err := r.unmarshalTraining(firestoreTraining)
		if err != nil {
			return err
		}

		if err := training.CanUserSeeTraining(user, *tr); err != nil {
			return err
		}

		updatedTraining, err := updateFn(ctx, tr)
		if err != nil {
			return err
		}

		return tx.Set(documentRef, r.marshalTraining(updatedTraining))
	})
}

func (r TrainingsFirestoreRepository) marshalTraining(tr *training.Training) TrainingModel {
	trainingModel := TrainingModel{
		UUID:     tr.UUID(),
		UserUUID: tr.UserUUID(),
		User:     tr.UserName(),
		Time:     tr.Time(),
		Notes:    tr.Notes(),
		Canceled: tr.IsCanceled(),
	}

	if tr.IsRescheduleProposed() {
		proposedBy := tr.MovedProposedBy().String()
		proposedTime := tr.ProposedNewTime()

		trainingModel.MoveProposedBy = &proposedBy
		trainingModel.ProposedTime = &proposedTime
	}

	return trainingModel
}

func (r TrainingsFirestoreRepository) unmarshalTraining(doc *firestore.DocumentSnapshot) (*training.Training, error) {
	trainingModel := TrainingModel{}
	err := doc.DataTo(&trainingModel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load document")
	}

	var moveProposedBy training.UserType
	if trainingModel.MoveProposedBy != nil {
		moveProposedBy, err = training.NewUserTypeFromString(*trainingModel.MoveProposedBy)
		if err != nil {
			return nil, err
		}
	}

	var proposedTime time.Time
	if trainingModel.ProposedTime != nil {
		proposedTime = *trainingModel.ProposedTime
	}

	return training.UnmarshalTrainingFromDatabase(
		trainingModel.UUID,
		trainingModel.UserUUID,
		trainingModel.User,
		trainingModel.Time,
		trainingModel.Notes,
		trainingModel.Canceled,
		proposedTime,
		moveProposedBy,
	)
}

func (r TrainingsFirestoreRepository) AllTrainings(ctx context.Context) ([]query.Training, error) {
	query := r.
		trainingsCollection().
		Query.
		Where("Time", ">=", time.Now().Add(-time.Hour*24)).
		Where("Canceled", "==", false)

	iter := query.Documents(ctx)

	return r.trainingModelsToQuery(iter)
}

func (r TrainingsFirestoreRepository) FindTrainingsForUser(ctx context.Context, userUUID string) ([]query.Training, error) {
	query := r.trainingsCollection().Query.
		Where("Time", ">=", time.Now().Add(-time.Hour*24)).
		Where("UserUuid", "==", userUUID).
		Where("Canceled", "==", false)

	iter := query.Documents(ctx)

	return r.trainingModelsToQuery(iter)
}

// warning: RemoveAllTrainings was designed for tests for doing data cleanups
func (r TrainingsFirestoreRepository) RemoveAllTrainings(ctx context.Context) error {
	for {
		iter := r.trainingsCollection().Limit(100).Documents(ctx)
		numDeleted := 0

		batch := r.firestoreClient.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return errors.Wrap(err, "unable to get document")
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to remove docs")
		}
	}
}

func (r TrainingsFirestoreRepository) trainingModelsToQuery(iter *firestore.DocumentIterator) ([]query.Training, error) {
	var trainings []query.Training

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		tr, err := r.unmarshalTraining(doc)
		if err != nil {
			return nil, err
		}

		queryTraining := query.Training{
			UUID:           tr.UUID(),
			UserUUID:       tr.UserUUID(),
			User:           tr.UserName(),
			Time:           tr.Time(),
			Notes:          tr.Notes(),
			CanBeCancelled: tr.CanBeCanceledForFree(),
		}

		if tr.IsRescheduleProposed() {
			proposedTime := tr.ProposedNewTime()
			queryTraining.ProposedTime = &proposedTime

			proposedBy := tr.MovedProposedBy().String()
			queryTraining.MoveProposedBy = &proposedBy
		}

		trainings = append(trainings, queryTraining)
	}

	sort.Slice(trainings, func(i, j int) bool { return trainings[i].Time.Before(trainings[j].Time) })

	return trainings, nil
}
