package hour_test

import (
	"testing"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/domain/hour"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHour_MakeNotAvailable(t *testing.T) {
	h, err := hour.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)

	require.NoError(t, h.MakeNotAvailable())
	assert.False(t, h.IsAvailable())
}

func TestHour_MakeNotAvailable_with_scheduled_training(t *testing.T) {
	h := newHourWithScheduledTraining(t)

	assert.Equal(t, hour.ErrTrainingScheduled, h.MakeNotAvailable())
}

func TestHour_MakeAvailable(t *testing.T) {
	h, err := hour.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)

	require.NoError(t, h.MakeNotAvailable())

	require.NoError(t, h.MakeAvailable())
	assert.True(t, h.IsAvailable())
}

func TestHour_MakeAvailable_with_scheduled_training(t *testing.T) {
	h := newHourWithScheduledTraining(t)

	assert.Equal(t, hour.ErrTrainingScheduled, h.MakeAvailable())
}

func TestHour_ScheduleTraining(t *testing.T) {
	h, err := hour.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)

	require.NoError(t, h.ScheduleTraining())

	assert.True(t, h.HasTrainingScheduled())
	assert.False(t, h.IsAvailable())
}

func TestHour_ScheduleTraining_with_not_available(t *testing.T) {
	h := newNotAvailableHour(t)
	assert.Equal(t, hour.ErrHourNotAvailable, h.ScheduleTraining())
}

func TestHour_CancelTraining(t *testing.T) {
	h := newHourWithScheduledTraining(t)

	require.NoError(t, h.CancelTraining())

	assert.False(t, h.HasTrainingScheduled())
	assert.True(t, h.IsAvailable())
}

func TestHour_CancelTraining_no_training_scheduled(t *testing.T) {
	h, err := hour.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)

	assert.Equal(t, hour.ErrNoTrainingScheduled, h.CancelTraining())
}

func newHourWithScheduledTraining(t *testing.T) *hour.Hour {
	h, err := hour.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)

	require.NoError(t, h.ScheduleTraining())

	return h
}

func newNotAvailableHour(t *testing.T) *hour.Hour {
	h, err := hour.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)

	require.NoError(t, h.MakeNotAvailable())

	return h
}
