package query

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/domain/hour"
)

type HourAvailabilityHandler struct {
	hourRepo hour.Repository
}

func NewHourAvailabilityHandler(hourRepo hour.Repository) HourAvailabilityHandler {
	if hourRepo == nil {
		panic("nil hourRepo")
	}

	return HourAvailabilityHandler{hourRepo: hourRepo}
}

func (h HourAvailabilityHandler) Handle(ctx context.Context, time time.Time) (bool, error) {
	hour, err := h.hourRepo.GetOrCreateHour(ctx, time)
	if err != nil {
		return false, err
	}

	return hour.IsAvailable(), nil
}
