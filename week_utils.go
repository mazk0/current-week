package main

import "time"

func getFirstAndLastDateOfWeek(year, week int) (time.Time, time.Time) {
	firstDateOfWeek := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	for firstDateOfWeek.Weekday() != time.Monday {
		firstDateOfWeek = firstDateOfWeek.AddDate(0, 0, 1)
	}
	firstDateOfWeek = firstDateOfWeek.AddDate(0, 0, (week-1)*7)
	lastDateOfWeek := firstDateOfWeek.AddDate(0, 0, 6)

	return firstDateOfWeek, lastDateOfWeek
}
