package adapters

import (
	"context"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/auth"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/app"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

type TrainingModel struct {
	UUID     string `firestore:"Uuid"`
	UserUUID string `firestore:"UserUuid"`
	User     string `firestore:"User"`

	Time  time.Time `firestore:"Time"`
	Notes string    `firestore:"Notes"`

	ProposedTime   *time.Time `firestore:"ProposedTime"`
	MoveProposedBy *string    `firestore:"MoveProposedBy"`
}

func (t TrainingModel) canBeCancelled() bool {
	return t.Time.Sub(time.Now()) > time.Hour*24
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

func (d TrainingsFirestoreRepository) trainingsCollection() *firestore.CollectionRef {
	return d.firestoreClient.Collection("trainings")
}

func (d TrainingsFirestoreRepository) AllTrainings(ctx context.Context) ([]app.Training, error) {
	query := d.trainingsCollection().Query.Where("Time", ">=", time.Now().Add(-time.Hour*24))

	iter := query.Documents(ctx)

	return trainingModelsToApp(iter)
}

func (d TrainingsFirestoreRepository) FindTrainingsForUser(ctx context.Context, user auth.User) ([]app.Training, error) {
	query := d.trainingsCollection().Query.
		Where("Time", ">=", time.Now().Add(-time.Hour*24)).
		Where("UserUuid", "==", user.UUID)

	iter := query.Documents(ctx)

	return trainingModelsToApp(iter)
}

func trainingModelsToApp(iter *firestore.DocumentIterator) ([]app.Training, error) {
	var trainings []app.Training

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		t := TrainingModel{}
		if err := doc.DataTo(&t); err != nil {
			return nil, err
		}

		trainings = append(trainings, app.Training(t))
	}

	sort.Slice(trainings, func(i, j int) bool { return trainings[i].Time.Before(trainings[j].Time) })

	return trainings, nil
}

func (d TrainingsFirestoreRepository) CreateTraining(ctx context.Context, training app.Training, createFn func() error) error {
	collection := d.trainingsCollection()

	trainingModel := TrainingModel(training)

	return d.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		docs, err := tx.Documents(collection.Where("Time", "==", trainingModel.Time)).GetAll()
		if err != nil {
			return errors.Wrap(err, "unable to get actual docs")
		}
		if len(docs) > 0 {
			return errors.Errorf("there is training already at %s", trainingModel.Time)
		}

		err = createFn()
		if err != nil {
			return err
		}

		return tx.Create(collection.Doc(trainingModel.UUID), trainingModel)
	})
}

func (d TrainingsFirestoreRepository) CancelTraining(ctx context.Context, trainingUUID string, deleteFn func(app.Training) error) error {
	trainingsCollection := d.trainingsCollection()

	return d.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		trainingDocumentRef := trainingsCollection.Doc(trainingUUID)

		firestoreTraining, err := tx.Get(trainingDocumentRef)
		if err != nil {
			return errors.Wrap(err, "unable to get actual docs")
		}

		training := TrainingModel{}
		err = firestoreTraining.DataTo(&training)
		if err != nil {
			return errors.Wrap(err, "unable to load document")
		}

		err = deleteFn(app.Training(training))
		if err != nil {
			return err
		}

		return tx.Delete(trainingDocumentRef)
	})
}

func (d TrainingsFirestoreRepository) RescheduleTraining(
	ctx context.Context,
	trainingUUID string,
	newTime time.Time,
	updateFn func(app.Training) (app.Training, error),
) error {
	collection := d.trainingsCollection()

	return d.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(d.trainingsCollection().Doc(trainingUUID))
		if err != nil {
			return errors.Wrap(err, "could not find training")
		}

		docs, err := tx.Documents(collection.Where("Time", "==", newTime)).GetAll()
		if err != nil {
			return errors.Wrap(err, "unable to get actual docs")
		}
		if len(docs) > 0 {
			return errors.Errorf("there is training already at %s", newTime)
		}

		var training TrainingModel
		err = doc.DataTo(&training)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal training")
		}

		updatedTraining, err := updateFn(app.Training(training))
		if err != nil {
			return err
		}

		return tx.Set(collection.Doc(training.UUID), TrainingModel(updatedTraining))
	})
}

func (d TrainingsFirestoreRepository) ApproveTrainingReschedule(ctx context.Context, trainingUUID string, updateFn func(app.Training) (app.Training, error)) error {
	return d.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(d.trainingsCollection().Doc(trainingUUID))
		if err != nil {
			return errors.Wrap(err, "could not find training")
		}

		var training TrainingModel
		err = doc.DataTo(&training)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal training")
		}

		updatedTraining, err := updateFn(app.Training(training))
		if err != nil {
			return err
		}

		return tx.Set(d.trainingsCollection().Doc(training.UUID), TrainingModel(updatedTraining))
	})
}

func (d TrainingsFirestoreRepository) RejectTrainingReschedule(ctx context.Context, trainingUUID string, updateFn func(app.Training) (app.Training, error)) error {
	return d.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(d.trainingsCollection().Doc(trainingUUID))
		if err != nil {
			return errors.Wrap(err, "could not find training")
		}

		var training TrainingModel
		err = doc.DataTo(&training)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal training")
		}

		updatedTraining, err := updateFn(app.Training(training))
		if err != nil {
			return err
		}

		return tx.Set(d.trainingsCollection().Doc(training.UUID), TrainingModel(updatedTraining))
	})
}
