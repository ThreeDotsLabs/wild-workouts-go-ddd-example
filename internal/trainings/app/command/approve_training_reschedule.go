package command

import (
	"context"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/logs"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
)

type ApproveTrainingReschedule struct {
	TrainingUUID string
	User         training.User
}

type ApproveTrainingRescheduleHandler struct {
	repo           training.Repository
	userService    UserService
	trainerService TrainerService
}

func NewApproveTrainingRescheduleHandler(
	repo training.Repository,
	userService UserService,
	trainerService TrainerService,
) ApproveTrainingRescheduleHandler {
	if repo == nil {
		panic("nil repo")
	}
	if userService == nil {
		panic("nil userService")
	}
	if trainerService == nil {
		panic("nil trainerService")
	}

	return ApproveTrainingRescheduleHandler{repo, userService, trainerService}
}

func (h ApproveTrainingRescheduleHandler) Handle(ctx context.Context, cmd ApproveTrainingReschedule) (err error) {
	defer func() {
		logs.LogCommandExecution("ApproveTrainingReschedule", cmd, err)
	}()

	return h.repo.UpdateTraining(
		ctx,
		cmd.TrainingUUID,
		cmd.User,
		func(ctx context.Context, tr *training.Training) (*training.Training, error) {
			originalTrainingTime := tr.Time()

			if err := tr.ApproveReschedule(cmd.User.Type()); err != nil {
				return nil, err
			}

			err := h.trainerService.MoveTraining(ctx, tr.Time(), originalTrainingTime)
			if err != nil {
				return nil, err
			}

			return tr, nil
		},
	)
}
