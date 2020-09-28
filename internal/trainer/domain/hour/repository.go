package hour

import (
	"context"
	"time"
)

type Repository interface {
	GetHour(ctx context.Context, hourTime time.Time) (*Hour, error)
	UpdateHour(
		ctx context.Context,
		hourTime time.Time,
		updateFn func(h *Hour) (*Hour, error),
	) error
}
