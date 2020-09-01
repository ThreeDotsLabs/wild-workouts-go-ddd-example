package adapters

import (
	"context"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/domain/hour"
)

type MemoryHourRepository struct {
	hours map[time.Time]hour.Hour
	lock  *sync.RWMutex

	hourFactory hour.Factory
}

func NewMemoryHourRepository(hourFactory hour.Factory) *MemoryHourRepository {
	if hourFactory.IsZero() {
		panic("missing hourFactory")
	}

	return &MemoryHourRepository{
		hours:       map[time.Time]hour.Hour{},
		lock:        &sync.RWMutex{},
		hourFactory: hourFactory,
	}
}

func (m MemoryHourRepository) GetOrCreateHour(_ context.Context, hourTime time.Time) (*hour.Hour, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.getOrCreateHour(hourTime)
}

func (m MemoryHourRepository) getOrCreateHour(hourTime time.Time) (*hour.Hour, error) {
	currentHour, ok := m.hours[hourTime]
	if !ok {
		return m.hourFactory.NewNotAvailableHour(hourTime)
	}

	// we don't store hours as pointers, but as values
	// thanks to that, we are sure that nobody can modify Hour without using UpdateHour
	return &currentHour, nil
}

func (m *MemoryHourRepository) UpdateHour(
	_ context.Context,
	hourTime time.Time,
	updateFn func(h *hour.Hour) (*hour.Hour, error),
) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	currentHour, err := m.getOrCreateHour(hourTime)
	if err != nil {
		return err
	}

	updatedHour, err := updateFn(currentHour)
	if err != nil {
		return err
	}

	m.hours[hourTime] = *updatedHour

	return nil
}
