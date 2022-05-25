package command

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/decorator"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/logs"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
	"github.com/sirupsen/logrus"
)

type RescheduleTraining struct {
	TrainingUUID string
	NewTime      time.Time

	User training.User

	NewNotes string
}

type RescheduleTrainingHandler decorator.CommandHandler[RescheduleTraining]

type rescheduleTrainingHandler struct {
	repo           training.Repository
	userService    UserService
	trainerService TrainerService
}

func NewRescheduleTrainingHandler(
	repo training.Repository,
	userService UserService,
	trainerService TrainerService,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) RescheduleTrainingHandler {
	if repo == nil {
		panic("nil repo")
	}
	if userService == nil {
		panic("nil userService")
	}
	if trainerService == nil {
		panic("nil trainerService")
	}

	return decorator.ApplyCommandDecorators[RescheduleTraining](
		rescheduleTrainingHandler{repo: repo, userService: userService, trainerService: trainerService},
		logger,
		metricsClient,
	)
}

func (h rescheduleTrainingHandler) Handle(ctx context.Context, cmd RescheduleTraining) (err error) {
	defer func() {
		logs.LogCommandExecution("RescheduleTraining", cmd, err)
	}()

	return h.repo.UpdateTraining(
		ctx,
		cmd.TrainingUUID,
		cmd.User,
		func(ctx context.Context, tr *training.Training) (*training.Training, error) {
			originalTrainingTime := tr.Time()

			if err := tr.UpdateNotes(cmd.NewNotes); err != nil {
				return nil, err
			}

			if err := tr.RescheduleTraining(cmd.NewTime); err != nil {
				return nil, err
			}

			err := h.trainerService.MoveTraining(ctx, cmd.NewTime, originalTrainingTime)
			if err != nil {
				return nil, err
			}

			return tr, nil
		},
	)
}
