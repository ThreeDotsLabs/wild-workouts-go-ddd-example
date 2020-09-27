package command

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/logs"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
	"github.com/pkg/errors"
)

type ScheduleTraining struct {
	TrainingUUID string

	UserUUID string
	UserName string

	TrainingTime time.Time
	Notes        string
}

type ScheduleTrainingHandler struct {
	repo           training.Repository
	userService    UserService
	trainerService TrainerService
}

func NewScheduleTrainingHandler(repo training.Repository, userService UserService, trainerService TrainerService) ScheduleTrainingHandler {
	if repo == nil {
		panic("nil repo")
	}
	if userService == nil {
		panic("nil repo")
	}
	if trainerService == nil {
		panic("nil trainerService")
	}

	return ScheduleTrainingHandler{repo: repo, userService: userService, trainerService: trainerService}
}

func (h ScheduleTrainingHandler) Handle(ctx context.Context, cmd ScheduleTraining) (err error) {
	defer func() {
		logs.LogCommandExecution("ScheduleTraining", cmd, err)
	}()

	tr, err := training.NewTraining(cmd.TrainingUUID, cmd.UserUUID, cmd.UserName, cmd.TrainingTime)
	if err != nil {
		return err
	}

	if err := h.repo.AddTraining(ctx, tr); err != nil {
		return err
	}

	err = h.userService.UpdateTrainingBalance(ctx, tr.UserUUID(), -1)
	if err != nil {
		return errors.Wrap(err, "unable to change trainings balance")
	}

	err = h.trainerService.ScheduleTraining(ctx, tr.Time())
	if err != nil {
		return errors.Wrap(err, "unable to schedule training")
	}

	return nil
}
