package main

import (
	"time"
)

const (
	minHour = 12
	maxHour = 20
)

// setDefaultAvailability adds missing hours to Date model if they were not set
func setDefaultAvailability(date DateModel) DateModel {

HoursLoop:
	for hour := minHour; hour <= maxHour; hour++ {
		hour := time.Date(date.Date.Year(), date.Date.Month(), date.Date.Day(), hour, 0, 0, 0, time.UTC)

		for i := range date.Hours {
			if date.Hours[i].Hour.Equal(hour) {
				continue HoursLoop
			}
		}
		newHour := HourModel{
			Available: false,
			Hour:      hour,
		}

		date.Hours = append(date.Hours, newHour)
	}

	return date
}

func addMissingDates(params *GetTrainerAvailableHoursParams, dates []DateModel) []DateModel {
	for day := params.DateFrom.UTC(); day.Before(params.DateTo) || day.Equal(params.DateTo); day = day.Add(time.Hour * 24) {
		found := false
		for _, date := range dates {
			if date.Date.Equal(day) {
				found = true
				break
			}
		}

		if !found {
			date := DateModel{
				Date: day,
			}
			date = setDefaultAvailability(date)
			dates = append(dates, date)
		}
	}

	return dates
}
