package command

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/logs"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainings/domain/training"
)

type RequestTrainingReschedule struct {
	TrainingUUID string
	NewTime      time.Time

	User training.User

	NewNotes string
}

type RequestTrainingRescheduleHandler struct {
	repo training.Repository
}

func NewRequestTrainingRescheduleHandler(repo training.Repository) RequestTrainingRescheduleHandler {
	if repo == nil {
		panic("nil repo service")
	}

	return RequestTrainingRescheduleHandler{repo: repo}
}

func (h RequestTrainingRescheduleHandler) Handle(ctx context.Context, cmd RequestTrainingReschedule) (err error) {
	defer func() {
		logs.LogCommandExecution("RequestTrainingReschedule", cmd, err)
	}()

	return h.repo.UpdateTraining(
		ctx,
		cmd.TrainingUUID,
		cmd.User,
		func(ctx context.Context, tr *training.Training) (*training.Training, error) {
			if err := tr.UpdateNotes(cmd.NewNotes); err != nil {
				return nil, err
			}

			tr.ProposeReschedule(cmd.NewTime, cmd.User.Type())

			return tr, nil
		},
	)
}
