package hour

import (
	"time"

	"github.com/pkg/errors"
)

type Hour struct {
	hour time.Time

	availability Availability
}

var (
	ErrNotFullHour    = errors.New("hour should be a full hour")
	ErrTooDistantDate = errors.Errorf("schedule can be only set for next %d weeks", MaxWeeksInTheFutureToSet)
	ErrPastHour       = errors.New("cannot create hour from past")
	ErrTooEarlyHour   = errors.Errorf("too early hour, min UTC hour: %d", MinUtcHour)
	ErrTooLateHour    = errors.Errorf("too late hour, max UTC hour: %d", MaxUtcHour)
)

const (
	// in theory it may be in some config, but let's dont overcomplicate, YAGNI!
	MaxWeeksInTheFutureToSet = 6
	MinUtcHour               = 12
	MaxUtcHour               = 20

	day  = time.Hour * 24
	week = day * 7
)

func NewAvailableHour(hour time.Time) (*Hour, error) {
	if err := validateTime(hour); err != nil {
		return nil, err
	}

	return &Hour{
		hour:         hour,
		availability: Available,
	}, nil
}

func NewNotAvailableHour(hour time.Time) (*Hour, error) {
	if err := validateTime(hour); err != nil {
		return nil, err
	}

	return &Hour{
		hour:         hour,
		availability: NotAvailable,
	}, nil
}

// UnmarshalHourFromRepository unmarshals Hour from the database.
//
// It should be used only for unmarshalling from the database!
// You can't use UnmarshalHourFromRepository as constructor - It may put domain into the invalid state!
func UnmarshalHourFromRepository(hour time.Time, availability Availability) (*Hour, error) {
	if err := validateTime(hour); err != nil {
		return nil, err
	}

	if availability.IsZero() {
		return nil, errors.New("empty availability")
	}

	return &Hour{
		hour:         hour,
		availability: availability,
	}, nil
}

func validateTime(hour time.Time) error {
	if !hour.Round(time.Hour).Equal(hour) {
		return ErrNotFullHour
	}

	if hour.After(time.Now().Add(week * MaxWeeksInTheFutureToSet)) {
		return ErrTooDistantDate
	}

	currentHour := time.Now().Truncate(time.Hour)
	if hour.Before(currentHour) || hour.Equal(currentHour) {
		return ErrPastHour
	}
	if hour.UTC().Hour() > MaxUtcHour {
		return ErrTooLateHour
	}
	if hour.UTC().Hour() < MinUtcHour {
		return ErrTooEarlyHour
	}

	return nil
}

func (h *Hour) Time() time.Time {
	return h.hour
}
