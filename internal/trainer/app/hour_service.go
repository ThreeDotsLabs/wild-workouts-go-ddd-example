package app

import (
	"context"
	"sort"
	"time"

	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/common/errors"
	"github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/internal/trainer/domain/hour"
)

type HourService struct {
	datesRepo datesRepository
	hourRepo  hour.Repository
}

type datesRepository interface {
	GetDates(ctx context.Context, from time.Time, to time.Time) ([]Date, error)
}

func NewHourService(datesRepo datesRepository, hourRepo hour.Repository) HourService {
	return HourService{
		datesRepo: datesRepo,
		hourRepo:  hourRepo,
	}
}

func (c HourService) GetTrainerAvailableHours(ctx context.Context, from time.Time, to time.Time) ([]Date, error) {
	if from.After(to) {
		return nil, errors.NewIncorrectInputError("date-from-after-date-to", "Date from after date to")
	}

	dates, err := c.datesRepo.GetDates(ctx, from, to)
	if err != nil {
		return nil, err
	}

	dates = addMissingDates(dates, from, to)

	for i, date := range dates {
		date = setDefaultAvailability(date)
		sort.Slice(date.Hours, func(i, j int) bool { return date.Hours[i].Hour.Before(date.Hours[j].Hour) })
		dates[i] = date
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Date.Before(dates[j].Date) })

	return dates, nil
}

func (c HourService) MakeHoursAvailable(ctx context.Context, hours []time.Time) error {
	for _, hourToUpdate := range hours {
		if err := c.hourRepo.UpdateHour(ctx, hourToUpdate, func(h *hour.Hour) (*hour.Hour, error) {
			if err := h.MakeAvailable(); err != nil {
				return nil, err
			}
			return h, nil
		}); err != nil {
			return errors.NewSlugError(err.Error(), "unable-to-update-availability")
		}
	}

	return nil
}

func (c HourService) MakeHoursUnavailable(ctx context.Context, hours []time.Time) error {
	for _, hourToUpdate := range hours {
		if err := c.hourRepo.UpdateHour(ctx, hourToUpdate, func(h *hour.Hour) (*hour.Hour, error) {
			if err := h.MakeNotAvailable(); err != nil {
				return nil, err
			}
			return h, nil
		}); err != nil {
			return errors.NewSlugError(err.Error(), "unable-to-update-availability")
		}
	}

	return nil
}
