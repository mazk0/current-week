package main

import (
	"net/http"
	"net/url"
	"testing"
)

func TestParseWeekYearFromRequest(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantYear    int
		wantWeek    int
		expectError bool
	}{
		{
			name:     "Valid year and week",
			path:     "/year/2022/week/10",
			wantYear: 2022,
			wantWeek: 10,
		},
		{
			name:     "Valid year only",
			path:     "/year/2023",
			wantYear: 2023,
			wantWeek: -1,
		},
		{
			name:     "Valid week only",
			path:     "/week/42",
			wantYear: -1,
			wantWeek: 42,
		},
		{
			name:        "Invalid year",
			path:        "/year/abc",
			expectError: true,
		},
		{
			name:        "Invalid week",
			path:        "/week/xyz",
			expectError: true,
		},
		{
			name:     "No params",
			path:     "/",
			wantYear: -1,
			wantWeek: -1,
		},
		{
			name:     "Partial match year",
			path:     "/year/",
			wantYear: -1,
			wantWeek: -1,
		},
		{
			name:     "Param without value at end",
			path:     "/year",
			wantYear: -1,
			wantWeek: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				URL: &url.URL{Path: tt.path},
			}

			year, week, err := parseWeekYearFromRequest(req)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.wantYear != -1 && year != tt.wantYear {
				t.Errorf("Year: got %d, want %d", year, tt.wantYear)
			}
			if tt.wantWeek != -1 && week != tt.wantWeek {
				t.Errorf("Week: got %d, want %d", week, tt.wantWeek)
			}
		})
	}
}

func BenchmarkParseWeekYearFromRequest(b *testing.B) {
	req := &http.Request{
		URL: &url.URL{Path: "/year/2023/week/42"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseWeekYearFromRequest(req)
	}
}
