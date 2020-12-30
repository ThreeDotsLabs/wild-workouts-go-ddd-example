package hour_test

import (
	"testing"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/domain/hour"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testHourFactory = hour.MustNewFactory(hour.FactoryConfig{
	MaxWeeksInTheFutureToSet: 100,
	MinUtcHour:               0,
	MaxUtcHour:               24,
})

func TestNewAvailableHour(t *testing.T) {
	h, err := testHourFactory.NewAvailableHour(validTrainingHour())
	require.NoError(t, err)

	assert.True(t, h.IsAvailable())
}

func TestNewAvailableHour_not_full_hour(t *testing.T) {
	constructorTime := trainingHourWithMinutes(13)

	_, err := testHourFactory.NewAvailableHour(constructorTime)
	assert.Equal(t, hour.ErrNotFullHour, err)
}

func TestNewAvailableHour_too_distant_date(t *testing.T) {
	maxWeeksInFuture := 1

	factory := hour.MustNewFactory(hour.FactoryConfig{
		MaxWeeksInTheFutureToSet: maxWeeksInFuture,
		MinUtcHour:               0,
		MaxUtcHour:               0,
	})

	constructorTime := time.Now().Truncate(time.Hour*24).AddDate(0, 0, maxWeeksInFuture*7+1)

	_, err := factory.NewAvailableHour(constructorTime)
	assert.Equal(
		t,
		hour.TooDistantDateError{
			MaxWeeksInTheFutureToSet: maxWeeksInFuture,
			ProvidedDate:             constructorTime,
		},
		err,
	)
}

func TestNewAvailableHour_past_date(t *testing.T) {
	pastHour := time.Now().Truncate(time.Hour).Add(-time.Hour)
	_, err := testHourFactory.NewAvailableHour(pastHour)
	assert.Equal(t, hour.ErrPastHour, err)

	currentHour := time.Now().Truncate(time.Hour)
	_, err = testHourFactory.NewAvailableHour(currentHour)
	assert.Equal(t, hour.ErrPastHour, err)
}

func TestNewAvailableHour_too_early_hour(t *testing.T) {
	factory := hour.MustNewFactory(hour.FactoryConfig{
		MaxWeeksInTheFutureToSet: 10,
		MinUtcHour:               12,
		MaxUtcHour:               18,
	})

	// we are using next day, to be sure that provided hour is not in the past
	currentTime := time.Now().AddDate(0, 0, 1)

	tooEarlyHour := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		factory.Config().MinUtcHour-1, 0, 0, 0,
		time.UTC,
	)

	_, err := factory.NewAvailableHour(tooEarlyHour)
	assert.Equal(
		t,
		hour.TooEarlyHourError{
			MinUtcHour:   factory.Config().MinUtcHour,
			ProvidedTime: tooEarlyHour,
		},
		err,
	)
}

func TestNewAvailableHour_too_late_hour(t *testing.T) {
	factory := hour.MustNewFactory(hour.FactoryConfig{
		MaxWeeksInTheFutureToSet: 10,
		MinUtcHour:               12,
		MaxUtcHour:               18,
	})

	// we are using next day, to be sure that provided hour is not in the past
	currentTime := time.Now().AddDate(0, 0, 1)

	tooEarlyHour := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		factory.Config().MaxUtcHour+1, 0, 0, 0,
		time.UTC,
	)

	_, err := factory.NewAvailableHour(tooEarlyHour)
	assert.Equal(
		t,
		hour.TooLateHourError{
			MaxUtcHour:   factory.Config().MaxUtcHour,
			ProvidedTime: tooEarlyHour,
		},
		err,
	)
}

func TestHour_Time(t *testing.T) {
	expectedTime := validTrainingHour()

	h, err := testHourFactory.NewAvailableHour(expectedTime)
	require.NoError(t, err)

	assert.Equal(t, expectedTime, h.Time())
}

func TestUnmarshalHourFromDatabase(t *testing.T) {
	trainingTime := validTrainingHour()

	h, err := testHourFactory.UnmarshalHourFromDatabase(trainingTime, hour.TrainingScheduled)
	require.NoError(t, err)

	assert.Equal(t, trainingTime, h.Time())
	assert.True(t, h.HasTrainingScheduled())
}

func TestFactoryConfig_Validate(t *testing.T) {
	testCases := []struct {
		Name        string
		Config      hour.FactoryConfig
		ExpectedErr string
	}{
		{
			Name: "valid",
			Config: hour.FactoryConfig{
				MaxWeeksInTheFutureToSet: 10,
				MinUtcHour:               10,
				MaxUtcHour:               12,
			},
			ExpectedErr: "",
		},
		{
			Name: "equal_min_and_max_hour",
			Config: hour.FactoryConfig{
				MaxWeeksInTheFutureToSet: 10,
				MinUtcHour:               12,
				MaxUtcHour:               12,
			},
			ExpectedErr: "",
		},
		{
			Name: "min_hour_after_max_hour",
			Config: hour.FactoryConfig{
				MaxWeeksInTheFutureToSet: 10,
				MinUtcHour:               13,
				MaxUtcHour:               12,
			},
			ExpectedErr: "MinUtcHour (13) can't be after MaxUtcHour (12)",
		},
		{
			Name: "zero_max_weeks",
			Config: hour.FactoryConfig{
				MaxWeeksInTheFutureToSet: 0,
				MinUtcHour:               10,
				MaxUtcHour:               12,
			},
			ExpectedErr: "MaxWeeksInTheFutureToSet should be greater than 1, but is 0",
		},
		{
			Name: "sub_zero_min_hour",
			Config: hour.FactoryConfig{
				MaxWeeksInTheFutureToSet: 10,
				MinUtcHour:               -1,
				MaxUtcHour:               12,
			},
			ExpectedErr: "MinUtcHour should be value between 0 and 24, but is -1",
		},
		{
			Name: "sub_zero_max_hour",
			Config: hour.FactoryConfig{
				MaxWeeksInTheFutureToSet: 10,
				MinUtcHour:               10,
				MaxUtcHour:               -1,
			},
			ExpectedErr: "MinUtcHour should be value between 0 and 24, but is -1; MinUtcHour (10) can't be after MaxUtcHour (-1)",
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			err := c.Config.Validate()

			if c.ExpectedErr != "" {
				assert.EqualError(t, err, c.ExpectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewFactory_invalid_config(t *testing.T) {
	f, err := hour.NewFactory(hour.FactoryConfig{})
	assert.Error(t, err)
	assert.Zero(t, f)
}

func validTrainingHour() time.Time {
	tomorrow := time.Now().Add(time.Hour * 24)

	return time.Date(
		tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		testHourFactory.Config().MinUtcHour, 0, 0, 0,
		time.UTC,
	)
}

func trainingHourWithMinutes(minute int) time.Time {
	tomorrow := time.Now().Add(time.Hour * 24)

	return time.Date(
		tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		testHourFactory.Config().MaxUtcHour, minute, 0, 0,
		time.UTC,
	)
}
