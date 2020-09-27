package command

import (
	"context"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/logs"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
)

type RejectTrainingReschedule struct {
	TrainingUUID string
	User         training.User
}

type RejectTrainingRescheduleHandler struct {
	repo training.Repository
}

func NewRejectTrainingRescheduleHandler(repo training.Repository) RejectTrainingRescheduleHandler {
	if repo == nil {
		panic("nil repo service")
	}

	return RejectTrainingRescheduleHandler{repo: repo}
}

func (h RejectTrainingRescheduleHandler) Handle(ctx context.Context, cmd RejectTrainingReschedule) (err error) {
	defer func() {
		logs.LogCommandExecution("RejectTrainingReschedule", cmd, err)
	}()

	return h.repo.UpdateTraining(
		ctx,
		cmd.TrainingUUID,
		cmd.User,
		func(ctx context.Context, tr *training.Training) (*training.Training, error) {
			if err := tr.RejectReschedule(); err != nil {
				return nil, err
			}

			return tr, nil
		},
	)
}
