package main

import (
	"testing"
)

func TestGetFirstAndLastDateOfWeek(t *testing.T) {
	tests := []struct {
		year      int
		week      int
		wantStart string
		wantEnd   string
	}{
		{2023, 1, "2023-01-02", "2023-01-08"},
		{2023, 52, "2023-12-25", "2023-12-31"},
		{2024, 1, "2024-01-01", "2024-01-07"}, // Leap year starting on Monday
		{2024, 29, "2024-07-15", "2024-07-21"},
		{2025, 1, "2024-12-30", "2025-01-05"},
		{2020, 1, "2019-12-30", "2020-01-05"}, // Leap year
		{2020, 53, "2020-12-28", "2021-01-03"},
	}

	for _, tt := range tests {
		start, end := getFirstAndLastDateOfWeek(tt.year, tt.week)
		startStr := start.Format("2006-01-02")
		endStr := end.Format("2006-01-02")

		if startStr != tt.wantStart {
			t.Errorf("getFirstAndLastDateOfWeek(%d, %d) start = %s, want %s", tt.year, tt.week, startStr, tt.wantStart)
		}
		if endStr != tt.wantEnd {
			t.Errorf("getFirstAndLastDateOfWeek(%d, %d) end = %s, want %s", tt.year, tt.week, endStr, tt.wantEnd)
		}
	}
}

func BenchmarkGetFirstAndLastDateOfWeek(b *testing.B) {
	// Benchmark for a typical year
	year := 2024
	week := 25
	for i := 0; i < b.N; i++ {
		getFirstAndLastDateOfWeek(year, week)
	}
}
