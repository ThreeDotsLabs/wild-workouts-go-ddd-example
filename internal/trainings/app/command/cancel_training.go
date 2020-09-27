package command

import (
	"context"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/logs"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
	"github.com/pkg/errors"
)

type CancelTraining struct {
	TrainingUUID string
	User         training.User
}

type CancelTrainingHandler struct {
	repo           training.Repository
	userService    UserService
	trainerService TrainerService
}

func NewCancelTrainingHandler(repo training.Repository, userService UserService, trainerService TrainerService) CancelTrainingHandler {
	if repo == nil {
		panic("nil repo")
	}
	if userService == nil {
		panic("nil user service")
	}
	if trainerService == nil {
		panic("nil trainer service")
	}

	return CancelTrainingHandler{repo: repo, userService: userService, trainerService: trainerService}
}

func (h CancelTrainingHandler) Handle(ctx context.Context, cmd CancelTraining) (err error) {
	defer func() {
		logs.LogCommandExecution("CancelTrainingHandler", cmd, err)
	}()

	return h.repo.UpdateTraining(
		ctx,
		cmd.TrainingUUID,
		cmd.User,
		func(ctx context.Context, tr *training.Training) (*training.Training, error) {
			if err := tr.Cancel(); err != nil {
				return nil, err
			}

			if balanceDelta := training.CancelBalanceDelta(*tr, cmd.User.Type()); balanceDelta != 0 {
				err := h.userService.UpdateTrainingBalance(ctx, tr.UserUUID(), balanceDelta)
				if err != nil {
					return nil, errors.Wrap(err, "unable to change trainings balance")
				}
			}

			if err := h.trainerService.CancelTraining(ctx, tr.Time()); err != nil {
				return nil, errors.Wrap(err, "unable to cancel training")
			}

			return tr, nil
		},
	)
}
