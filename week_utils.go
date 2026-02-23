package main

import "time"

func getFirstAndLastDateOfWeek(year, week int) (time.Time, time.Time) {
	// January 4th is always in week 1 according to ISO 8601
	referenceDate := time.Date(year, time.January, 4, 0, 0, 0, 0, time.UTC)

	// Calculate the offset to get to the Monday of the first ISO week
	// Monday is 1, so if Weekday is Monday(1), offset is 0.
	// Sunday is 0, so if Weekday is Sunday(0), offset is 6.
	// (Weekday + 6) % 7 gives the correct offset (0 for Mon, 1 for Tue, ..., 6 for Sun)
	offset := (int(referenceDate.Weekday()) + 6) % 7

	// Move to the Monday of the first ISO week
	firstWeekMonday := referenceDate.AddDate(0, 0, -offset)

	// Now we move to the start of the desired week by adding (week - 1) weeks
	firstDateOfWeek := firstWeekMonday.AddDate(0, 0, (week-1)*7)
	lastDateOfWeek := firstDateOfWeek.AddDate(0, 0, 6)

	return firstDateOfWeek, lastDateOfWeek
}

func getNumberOfWeeks(year int) int {
	// December 28th is always in the last week according to ISO 8601
	endDate := time.Date(year, time.December, 28, 0, 0, 0, 0, time.UTC)
	_, lastWeek := endDate.ISOWeek()

	return lastWeek
}

func timeFromYearAndWeek(year, week int) time.Time {
	// January 4th is always in the first week of the year
	jan4 := time.Date(year, time.January, 4, 0, 0, 0, 0, time.UTC)
	// Find the Monday of the first week of the year
	_, isoWeek := jan4.ISOWeek()
	// Calculate the difference in weeks
	weekDiff := week - isoWeek
	// Add the difference in weeks to the Monday of the first week
	return jan4.AddDate(0, 0, weekDiff*7)
}
