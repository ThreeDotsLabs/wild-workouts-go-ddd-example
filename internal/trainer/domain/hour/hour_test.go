package hour_test

import (
	"testing"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/domain/hour"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	day  = time.Hour * 24
	week = day * 7
)

func TestNewAvailableHour(t *testing.T) {
	h, err := hour.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)

	assert.True(t, h.IsAvailable())
}

func TestNewAvailableHour_not_full_hour(t *testing.T) {
	constructorTime := trainingHourWithMinutes(13)

	_, err := hour.NewAvailableHour(constructorTime)
	assert.Equal(t, hour.ErrNotFullHour, err)
}

func TestNewAvailableHour_too_distant_date(t *testing.T) {
	constructorTime := time.Now().Truncate(day).Add(week * hour.MaxWeeksInTheFutureToSet).Add(day)

	_, err := hour.NewAvailableHour(constructorTime)
	assert.Equal(t, hour.ErrTooDistantDate, err)
}

func TestNewAvailableHour_past_date(t *testing.T) {
	pastHour := time.Now().Truncate(time.Hour).Add(-time.Hour)
	_, err := hour.NewAvailableHour(pastHour)
	assert.Equal(t, hour.ErrPastHour, err)

	currentHour := time.Now().Truncate(time.Hour)
	_, err = hour.NewAvailableHour(currentHour)
	assert.Equal(t, hour.ErrPastHour, err)
}

func TestNewAvailableHour_too_early_hour(t *testing.T) {
	currentTime := time.Now().Add(day)
	constructorTime := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		hour.MinUtcHour-1, 0, 0, 0,
		time.UTC,
	)

	_, err := hour.NewAvailableHour(constructorTime)
	assert.Equal(t, hour.ErrTooEarlyHour, err)
}

func TestNewAvailableHour_too_late_hour(t *testing.T) {
	currentTime := time.Now()
	constructorTime := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		hour.MaxUtcHour+1, 0, 0, 0,
		time.UTC,
	)

	_, err := hour.NewAvailableHour(constructorTime)
	assert.Equal(t, hour.ErrTooLateHour, err)
}

func TestHour_Time(t *testing.T) {
	expectedTime := validTrainingHour()

	h, err := hour.NewAvailableHour(expectedTime)
	require.NoError(t, err)

	assert.Equal(t, expectedTime, h.Time())
}

func TestUnmarshalHourFromRepository(t *testing.T) {
	trainingTime := validTrainingHour()

	h, err := hour.UnmarshalHourFromRepository(trainingTime, hour.TrainingScheduled)
	require.NoError(t, err)

	assert.Equal(t, trainingTime, h.Time())
	assert.True(t, h.HasTrainingScheduled())
}

func validTrainingHour() time.Time {
	tomorrow := time.Now().Add(day)

	return time.Date(
		tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		hour.MinUtcHour, 0, 0, 0,
		time.UTC,
	)
}

func trainingHourWithMinutes(minute int) time.Time {
	tomorrow := time.Now().Add(day)

	return time.Date(
		tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		hour.MinUtcHour, minute, 0, 0,
		time.UTC,
	)
}
