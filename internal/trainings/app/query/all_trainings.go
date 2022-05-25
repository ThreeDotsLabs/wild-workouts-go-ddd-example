package query

import (
	"context"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/decorator"
	"github.com/sirupsen/logrus"
)

type AllTrainings struct{}

type AllTrainingsHandler decorator.QueryHandler[AllTrainings, []Training]

type allTrainingsHandler struct {
	readModel AllTrainingsReadModel
}

func NewAllTrainingsHandler(
	readModel AllTrainingsReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) AllTrainingsHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[AllTrainings, []Training](
		allTrainingsHandler{readModel: readModel},
		logger,
		metricsClient,
	)
}

type AllTrainingsReadModel interface {
	AllTrainings(ctx context.Context) ([]Training, error)
}

func (h allTrainingsHandler) Handle(ctx context.Context, _ AllTrainings) (tr []Training, err error) {
	return h.readModel.AllTrainings(ctx)
}
