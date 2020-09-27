package query

import (
	"context"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/auth"
)

type TrainingsForUserHandler struct {
	readModel TrainingsForUserReadModel
}

func NewTrainingsForUserHandler(readModel TrainingsForUserReadModel) TrainingsForUserHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return TrainingsForUserHandler{readModel: readModel}
}

type TrainingsForUserReadModel interface {
	FindTrainingsForUser(ctx context.Context, userUUID string) ([]Training, error)
}

func (h TrainingsForUserHandler) Handle(ctx context.Context, user auth.User) (tr []Training, err error) {
	return h.readModel.FindTrainingsForUser(ctx, user.UUID)
}
