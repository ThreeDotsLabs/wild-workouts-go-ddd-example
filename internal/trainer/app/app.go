package app

import (
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/app/command"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CancelTraining   command.CancelTrainingHandler
	ScheduleTraining command.ScheduleTrainingHandler

	MakeHoursAvailable   command.MakeHoursAvailableHandler
	MakeHoursUnavailable command.MakeHoursUnavailableHandler
}

type Queries struct {
	HourAvailability      query.HourAvailabilityHandler
	TrainerAvailableHours query.AvailableHoursHandler
}
