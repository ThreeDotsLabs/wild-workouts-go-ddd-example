package main

import (
	"context"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/auth"
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
