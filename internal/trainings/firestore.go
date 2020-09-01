package main

import (
	"context"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/auth"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/genproto/trainer"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/genproto/users"
	"github.com/golang/protobuf/ptypes"
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

type db struct {
	firestoreClient *firestore.Client
	trainerClient   trainer.TrainerServiceClient
	usersClient     users.UsersServiceClient
}

func (d db) TrainingsCollection() *firestore.CollectionRef {
	return d.firestoreClient.Collection("trainings")
}

func (d db) GetTrainings(ctx context.Context, user auth.User) ([]TrainingModel, error) {
	query := d.TrainingsCollection().Query.Where("Time", ">=", time.Now().Add(-time.Hour*24))

	if user.Role != "trainer" {
		query = query.Where("UserUuid", "==", user.UUID)
	}

	iter := query.Documents(ctx)

	var trainings []TrainingModel

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

		trainings = append(trainings, t)
	}

	sort.Slice(trainings, func(i, j int) bool { return trainings[i].Time.Before(trainings[j].Time) })

	return trainings, nil
}

func (d db) CreateTraining(ctx context.Context, user auth.User, training TrainingModel) error {
	collection := d.TrainingsCollection()

	return d.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		docs, err := tx.Documents(collection.Where("Time", "==", training.Time)).GetAll()
		if err != nil {
			return errors.Wrap(err, "unable to get actual docs")
		}
		if len(docs) > 0 {
			return errors.Errorf("there is training already at %s", training.Time)
		}

		_, err = d.usersClient.UpdateTrainingBalance(ctx, &users.UpdateTrainingBalanceRequest{
			UserId:       user.UUID,
			AmountChange: -1,
		})
		if err != nil {
			return errors.Wrap(err, "unable to change trainings balance")
		}

		timestamp, err := ptypes.TimestampProto(training.Time)
		if err != nil {
			return errors.Wrap(err, "unable to convert time to proto timestamp")
		}
		_, err = d.trainerClient.ScheduleTraining(ctx, &trainer.UpdateHourRequest{
			Time: timestamp,
		})
		if err != nil {
			return errors.Wrap(err, "unable to update trainer hour")
		}

		return tx.Create(collection.Doc(training.UUID), training)
	})
}

func (d db) CancelTraining(ctx context.Context, user auth.User, trainingUUID string) error {
	trainingsCollection := d.TrainingsCollection()

	return d.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		trainingDocumentRef := trainingsCollection.Doc(trainingUUID)

		firestoreTraining, err := tx.Get(trainingDocumentRef)
		if err != nil {
			return errors.Wrap(err, "unable to get actual docs")
		}

		training := &TrainingModel{}
		err = firestoreTraining.DataTo(training)
		if err != nil {
			return errors.Wrap(err, "unable to load document")
		}

		if user.Role != "trainer" && training.UserUUID != user.UUID {
			return errors.Errorf("user '%s' is trying to cancel training of user '%s'", user.UUID, training.UserUUID)
		}

		var trainingBalanceDelta int64
		if training.canBeCancelled() {
			// just give training back
			trainingBalanceDelta = 1
		} else {
			if user.Role == "trainer" {
				// 1 for cancelled training +1 fine for cancelling by trainer less than 24h before training
				trainingBalanceDelta = 2
			} else {
				// fine for cancelling less than 24h before training
				trainingBalanceDelta = 0
			}
		}

		if trainingBalanceDelta != 0 {
			_, err := d.usersClient.UpdateTrainingBalance(ctx, &users.UpdateTrainingBalanceRequest{
				UserId:       training.UserUUID,
				AmountChange: trainingBalanceDelta,
			})
			if err != nil {
				return errors.Wrap(err, "unable to change trainings balance")
			}
		}

		timestamp, err := ptypes.TimestampProto(training.Time)
		if err != nil {
			return errors.Wrap(err, "unable to convert time to proto timestamp")
		}
		_, err = d.trainerClient.CancelTraining(ctx, &trainer.UpdateHourRequest{
			Time: timestamp,
		})
		if err != nil {
			return errors.Wrap(err, "unable to update trainer hour")
		}

		return tx.Delete(trainingDocumentRef)
	})
}

func (d db) RescheduleTraining(ctx context.Context, user auth.User, trainingUUID string, newTime time.Time, notes string) error {
	collection := d.TrainingsCollection()

	return d.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(d.TrainingsCollection().Doc(trainingUUID))
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

		if training.canBeCancelled() {
			err = d.rescheduleTraining(ctx, training.Time, newTime)
			if err != nil {
				return errors.Wrap(err, "unable to reschedule training")
			}

			training.Time = newTime
			training.Notes = notes
		} else {
			training.ProposedTime = &newTime
			training.MoveProposedBy = &user.Role
			training.Notes = notes
		}

		return tx.Set(collection.Doc(training.UUID), training)
	})
}
func (d db) rescheduleTraining(ctx context.Context, oldTime, newTime time.Time) error {
	oldTimeProto, err := ptypes.TimestampProto(oldTime)
	if err != nil {
		return errors.Wrap(err, "unable to convert time to proto timestamp")
	}

	newTimeProto, err := ptypes.TimestampProto(newTime)
	if err != nil {
		return errors.Wrap(err, "unable to convert time to proto timestamp")
	}

	_, err = d.trainerClient.ScheduleTraining(ctx, &trainer.UpdateHourRequest{
		Time: newTimeProto,
	})
	if err != nil {
		return errors.Wrap(err, "unable to update trainer hour")
	}

	_, err = d.trainerClient.CancelTraining(ctx, &trainer.UpdateHourRequest{
		Time: oldTimeProto,
	})
	if err != nil {
		return errors.Wrap(err, "unable to update trainer hour")
	}

	return nil
}

func (d db) ApproveTrainingReschedule(ctx context.Context, user auth.User, trainingUUID string) error {
	return d.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(d.TrainingsCollection().Doc(trainingUUID))
		if err != nil {
			return errors.Wrap(err, "could not find training")
		}

		var training TrainingModel
		err = doc.DataTo(&training)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal training")
		}

		if training.ProposedTime == nil {
			return errors.New("training has no proposed time")
		}
		if training.MoveProposedBy == nil {
			return errors.New("training has no MoveProposedBy")
		}
		if *training.MoveProposedBy == "trainer" && training.UserUUID != user.UUID {
			return errors.Errorf("user '%s' cannot approve reschedule of user '%s'", user.UUID, training.UserUUID)
		}
		if *training.MoveProposedBy == user.Role {
			return errors.New("reschedule cannot be accepted by requesting person")
		}

		training.Time = *training.ProposedTime
		training.ProposedTime = nil

		return tx.Set(d.TrainingsCollection().Doc(training.UUID), training)
	})
}

func (d db) RejectTrainingReschedule(ctx context.Context, user auth.User, trainingUUID string) error {
	return d.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(d.TrainingsCollection().Doc(trainingUUID))
		if err != nil {
			return errors.Wrap(err, "could not find training")
		}

		var training TrainingModel
		err = doc.DataTo(&training)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal training")
		}

		if training.MoveProposedBy == nil {
			return errors.New("training has no MoveProposedBy")
		}
		if *training.MoveProposedBy != "trainer" && training.UserUUID != user.UUID {
			return errors.Errorf("user '%s' cannot approve reschedule of user '%s'", user.UUID, training.UserUUID)
		}

		training.ProposedTime = nil

		return tx.Set(d.TrainingsCollection().Doc(training.UUID), training)
	})
}
