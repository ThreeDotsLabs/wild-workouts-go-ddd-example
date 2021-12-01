package adapters

import (
	"context"
	"sort"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/app/query"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/domain/hour"
	"github.com/pkg/errors"

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
	factoryConfig   hour.FactoryConfig
}

func NewDatesFirestoreRepository(firestoreClient *firestore.Client, factoryConfig hour.FactoryConfig) DatesFirestoreRepository {
	if firestoreClient == nil {
		panic("missing firestoreClient")
	}

	return DatesFirestoreRepository{firestoreClient, factoryConfig}
}

func (d DatesFirestoreRepository) trainerHoursCollection() *firestore.CollectionRef {
	return d.firestoreClient.Collection("trainer-hours")
}

func (d DatesFirestoreRepository) DocumentRef(dateTimeToUpdate time.Time) *firestore.DocumentRef {
	return d.trainerHoursCollection().Doc(dateTimeToUpdate.Format("2006-01-02"))
}

func (d DatesFirestoreRepository) AvailableHours(ctx context.Context, from time.Time, to time.Time) ([]query.Date, error) {
	iter := d.
		trainerHoursCollection().
		Where("Date", ">=", from).
		Where("Date", "<=", to).
		Documents(ctx)

	var dates []query.Date

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
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

	dates = addMissingDates(dates, from, to)
	for i, date := range dates {
		date = d.setDefaultAvailability(date)
		sort.Slice(date.Hours, func(i, j int) bool { return date.Hours[i].Hour.Before(date.Hours[j].Hour) })
		dates[i] = date
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Date.Before(dates[j].Date) })

	return dates, nil
}

// setDefaultAvailability adds missing hours to Date model if they were not set
func (d DatesFirestoreRepository) setDefaultAvailability(date query.Date) query.Date {
HoursLoop:
	for h := d.factoryConfig.MinUtcHour; h <= d.factoryConfig.MaxUtcHour; h++ {
		hour := time.Date(date.Date.Year(), date.Date.Month(), date.Date.Day(), h, 0, 0, 0, time.UTC)

		for i := range date.Hours {
			if date.Hours[i].Hour.Equal(hour) {
				continue HoursLoop
			}
		}
		newHour := query.Hour{
			Available: false,
			Hour:      hour,
		}

		date.Hours = append(date.Hours, newHour)
	}

	return date
}

func addMissingDates(dates []query.Date, from time.Time, to time.Time) []query.Date {
	for day := from.UTC(); day.Before(to) || day.Equal(to); day = day.AddDate(0, 0, 1) {
		found := false
		for _, date := range dates {
			if date.Date.Equal(day) {
				found = true
				break
			}
		}

		if !found {
			date := query.Date{
				Date: day,
			}
			dates = append(dates, date)
		}
	}

	return dates
}

func dateModelToApp(dm DateModel) query.Date {
	var hours []query.Hour
	for _, h := range dm.Hours {
		hours = append(hours, query.Hour{
			Available:            h.Available,
			HasTrainingScheduled: h.HasTrainingScheduled,
			Hour:                 h.Hour,
		})
	}

	return query.Date{
		Date:         dm.Date,
		HasFreeHours: dm.HasFreeHours,
		Hours:        hours,
	}
}
