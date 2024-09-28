package main

import "time"

func getFirstAndLastDateOfWeek(year, week int) (time.Time, time.Time) {
	// January 4th is always in week 1 according to ISO 8601
	referenceDate := time.Date(year, time.January, 4, 0, 0, 0, 0, time.UTC)

	// Find the first ISO week of the year
	isoYear, isoWeek := referenceDate.ISOWeek()
	for isoYear < year || (isoYear == year && isoWeek != 1) {
		referenceDate = referenceDate.AddDate(0, 0, -1)
		isoYear, isoWeek = referenceDate.ISOWeek()
	}

	// Move back to the Monday of the first ISO week
	for referenceDate.Weekday() != time.Monday {
		referenceDate = referenceDate.AddDate(0, 0, -1)
	}

	// Now we move to the start of the desired week by adding (week - 1) weeks
	firstDateOfWeek := referenceDate.AddDate(0, 0, (week-1)*7)
	lastDateOfWeek := firstDateOfWeek.AddDate(0, 0, 6)

	return firstDateOfWeek, lastDateOfWeek
}

func getNumberOfWeeks(year int) int {
	// December 28th is always in the last week according to ISO 8601
	endDate := time.Date(year, time.December, 28, 0, 0, 0, 0, time.UTC)
	_, lastWeek := endDate.ISOWeek()

	return lastWeek
}
