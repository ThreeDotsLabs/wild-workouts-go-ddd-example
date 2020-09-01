package main

import (
	"context"
	"sort"
	"time"

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

type db struct {
	firestoreClient *firestore.Client
}

func (d db) TrainerHoursCollection() *firestore.CollectionRef {
	return d.firestoreClient.Collection("trainer-hours")
}

func (d db) DocumentRef(dateTimeToUpdate time.Time) *firestore.DocumentRef {
	return d.TrainerHoursCollection().Doc(dateTimeToUpdate.Format("2006-01-02"))
}

func (d db) GetDates(ctx context.Context, params *GetTrainerAvailableHoursParams) ([]DateModel, error) {
	dates, err := d.QueryDates(params, ctx)
	if err != nil {
		return nil, err
	}
	dates = addMissingDates(params, dates)

	for _, date := range dates {
		sort.Slice(date.Hours, func(i, j int) bool { return date.Hours[i].Hour.Before(date.Hours[j].Hour) })
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Date.Before(dates[j].Date) })

	return dates, nil
}

func (d db) QueryDates(params *GetTrainerAvailableHoursParams, ctx context.Context) ([]DateModel, error) {
	iter := d.
		TrainerHoursCollection().
		Where("Date", ">=", params.DateFrom).
		Where("Date", "<=", params.DateTo).
		Documents(ctx)

	var dates []DateModel

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
		date = setDefaultAvailability(date)
		dates = append(dates, date)
	}

	return dates, nil
}
