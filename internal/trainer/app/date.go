package app

import (
	"time"
)

type Date struct {
	Date         time.Time
	HasFreeHours bool
	Hours        []Hour
}

type Hour struct {
	Available            bool
	HasTrainingScheduled bool
	Hour                 time.Time
}

func (d Date) FindHourInDate(timeToCheck time.Time) (*Hour, bool) {
	for i, hour := range d.Hours {
		if hour.Hour == timeToCheck {
			return &d.Hours[i], true
		}
	}

	return nil, false
}

type AvailableHoursRequest struct {
	DateFrom time.Time
	DateTo   time.Time
}

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

func addMissingDates(dates []Date, from time.Time, to time.Time) []Date {
	for day := from.UTC(); day.Before(to) || day.Equal(to); day = day.Add(time.Hour * 24) {
		found := false
		for _, date := range dates {
			if date.Date.Equal(day) {
				found = true
				break
			}
		}

		if !found {
			date := Date{
				Date: day,
			}
			date = setDefaultAvailability(date)
			dates = append(dates, date)
		}
	}

	return dates
}
