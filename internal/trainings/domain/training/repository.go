package training

import (
	"context"
	"fmt"
)

type NotFoundError struct {
	TrainingUUID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("training '%s' not found", e.TrainingUUID)
}

type Repository interface {
	AddTraining(ctx context.Context, tr *Training) error

	GetTraining(ctx context.Context, trainingUUID string, user User) (*Training, error)

	UpdateTraining(
		ctx context.Context,
		trainingUUID string,
		user User,
		updateFn func(ctx context.Context, tr *Training) (*Training, error),
	) error
}
