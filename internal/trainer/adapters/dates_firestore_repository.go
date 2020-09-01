package adapters

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/app"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type DateModel struct {
	Date         time.Time   `firestore:"Date"`
	HasFreeHours bool        `firestore:"HasFreeHours"`
	Hours        []HourModel `firestore:"Hours"`
}

type HourModel struct {
	Available            bool      `firestore:"Available"`
	HasTrainingScheduled bool      `firestore:"HasTrainingScheduled"`
	Hour                 time.Time `firestore:"Hour"`
}

type DatesFirestoreRepository struct {
	firestoreClient *firestore.Client
}

func NewDatesFirestoreRepository(firestoreClient *firestore.Client) DatesFirestoreRepository {
	if firestoreClient == nil {
		panic("missing firestoreClient")
	}

	return DatesFirestoreRepository{
		firestoreClient: firestoreClient,
	}
}

func (d DatesFirestoreRepository) trainerHoursCollection() *firestore.CollectionRef {
	return d.firestoreClient.Collection("trainer-hours")
}

func (d DatesFirestoreRepository) DocumentRef(dateTimeToUpdate time.Time) *firestore.DocumentRef {
	return d.trainerHoursCollection().Doc(dateTimeToUpdate.Format("2006-01-02"))
}

func (d DatesFirestoreRepository) GetDates(ctx context.Context, from time.Time, to time.Time) ([]app.Date, error) {
	iter := d.
		trainerHoursCollection().
		Where("Date", ">=", from).
		Where("Date", "<=", to).
		Documents(ctx)

	var dates []app.Date

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		date := DateModel{}
		if err := doc.DataTo(&date); err != nil {
			return nil, err
		}
		dates = append(dates, dateModelToApp(date))
	}

	return dates, nil
}

func dateModelToApp(dm DateModel) app.Date {
	var hours []app.Hour
	for _, h := range dm.Hours {
		hours = append(hours, app.Hour{
			Available:            h.Available,
			HasTrainingScheduled: h.HasTrainingScheduled,
			Hour:                 h.Hour,
		})
	}

	return app.Date{
		Date:         dm.Date,
		HasFreeHours: dm.HasFreeHours,
		Hours:        hours,
	}
}

func (d DatesFirestoreRepository) CanLoadFixtures(ctx context.Context, daysToSet int) (bool, error) {
	documents, err := d.trainerHoursCollection().Limit(daysToSet).Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}

	return len(documents) < daysToSet, nil
}
