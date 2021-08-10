package service

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/adapters"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/app"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/app/command"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/app/query"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/domain/hour"
)

func NewApplication(ctx context.Context) app.Application {
	firestoreClient := &firestore.Client{}
	//if err != nil {
	//	panic(err)
	//}

	factoryConfig := hour.FactoryConfig{
		MaxWeeksInTheFutureToSet: 6,
		MinUtcHour:               12,
		MaxUtcHour:               20,
	}

	datesRepository := adapters.NewDatesFirestoreRepository(firestoreClient, factoryConfig)

	hourFactory, err := hour.NewFactory(factoryConfig)
	if err != nil {
		panic(err)
	}

	hourRepository := adapters.NewFirestoreHourRepository(firestoreClient, hourFactory)

	return app.Application{
		Commands: app.Commands{
			CancelTraining:       command.NewCancelTrainingHandler(hourRepository),
			ScheduleTraining:     command.NewScheduleTrainingHandler(hourRepository),
			MakeHoursAvailable:   command.NewMakeHoursAvailableHandler(hourRepository),
			MakeHoursUnavailable: command.NewMakeHoursUnavailableHandler(hourRepository),
		},
		Queries: app.Queries{
			HourAvailability:      query.NewHourAvailabilityHandler(hourRepository),
			TrainerAvailableHours: query.NewAvailableHoursHandler(datesRepository),
		},
	}
}
