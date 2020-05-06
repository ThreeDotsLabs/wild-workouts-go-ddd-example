package main

import (
	"context"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/auth"
	"google.golang.org/api/iterator"
)

type db struct {
	firestoreClient *firestore.Client
}

func (d db) TrainingsCollection() *firestore.CollectionRef {
	return d.firestoreClient.Collection("trainings")
}

func (d db) GetTrainings(ctx context.Context, user auth.User) ([]Training, error) {
	query := d.TrainingsCollection().Query.Where("Time", ">=", time.Now().Add(-time.Hour*24))

	if user.Role != "trainer" {
		query = query.Where("UserUuid", "==", user.UUID)
	}

	iter := query.Documents(ctx)

	var trainings []Training

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		t := Training{}
		if err := doc.DataTo(&t); err != nil {
			return nil, err
		}

		t.CanBeCancelled = t.canBeCancelled()
		t.MoveRequiresAccept = !t.canBeCancelled()

		trainings = append(trainings, t)
	}

	sort.Slice(trainings, func(i, j int) bool { return trainings[i].Time.Before(trainings[j].Time) })

	return trainings, nil
}
