package main

import (
	"time"

	"github.com/deepmap/oapi-codegen/pkg/types"
)

const (
	minHour = 12
	maxHour = 20
)

// setDefaultAvailability adds missing hours to Date model if they were not set
func setDefaultAvailability(date Date) Date {

HoursLoop:
	for hour := minHour; hour <= maxHour; hour++ {
		hour := time.Date(date.Date.Year(), date.Date.Month(), date.Date.Day(), hour, 0, 0, 0, time.UTC)

		for i := range date.Hours {
			if date.Hours[i].Hour.Equal(hour) {
				continue HoursLoop
			}
		}
		newHour := Hour{
			Available: false,
			Hour:      hour,
		}

		date.Hours = append(date.Hours, newHour)
	}

	return date
}

func addMissingDates(params *GetTrainerAvailableHoursParams, dates []Date) []Date {
	for day := params.DateFrom.UTC(); day.Before(params.DateTo) || day.Equal(params.DateTo); day = day.Add(time.Hour * 24) {
		found := false
		for _, date := range dates {
			if date.Date.Equal(day) {
				found = true
				break
			}
		}

		if !found {
			date := Date{
				Date: types.Date{
					Time: day,
				},
			}
			date = setDefaultAvailability(date)
			dates = append(dates, date)
		}
	}

	return dates
}

func (d Date) FindHourInDate(timeToCheck time.Time) (*Hour, bool) {
	for i, hour := range d.Hours {
		if hour.Hour == timeToCheck {
			return &d.Hours[i], true
		}
	}

	return nil, false
}
