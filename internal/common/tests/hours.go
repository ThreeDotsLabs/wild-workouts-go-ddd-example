package tests

import "time"

// RelativeDate returns a date in the future in specified days and hours.
// This allows running the tests in parallel, since each tests uses different date and hour.
//
// The downside of this approach is that you need to be aware of used dates when adding a new test.
// In our case this is not an issue, as it's trivial to see all usages and there's just a few of them.
//
// Another, more complex approach, would be to use random dates and retry in case of an error.
func RelativeDate(days int, hour int) time.Time {
	now := time.Now().UTC().AddDate(0, 0, days)
	return time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, time.UTC)
}
