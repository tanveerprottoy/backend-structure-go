package timeext

import (
	"testing"
	"time"
)

func TestNowUnixMilli(t *testing.T) {
	before := time.Now().UnixMilli()
	result := NowUnixMilli()
	after := time.Now().UnixMilli()

	if result < before || result > after {
		t.Errorf("NowUnixMilli() returned %d, expected value between %d and %d", result, before, after)
	}
}

func TestAddDate(t *testing.T) {
	tests := []struct {
		name     string
		years    int
		months   int
		days     int
		validate func(time.Time) bool
	}{
		{
			name:   "add 1 year",
			years:  1,
			months: 0,
			days:   0,
			validate: func(result time.Time) bool {
				return result.Year() == time.Now().Year()+1
			},
		},
		{
			name:   "add 1 month",
			years:  0,
			months: 1,
			days:   0,
			validate: func(result time.Time) bool {
				expectedMonth := time.Now().Month() + 1
				if expectedMonth > 12 {
					return result.Month() == expectedMonth-12 && result.Year() == time.Now().Year()+1
				}
				return result.Month() == expectedMonth && result.Year() == time.Now().Year()
			},
		},
		{
			name:   "add 1 day",
			years:  0,
			months: 0,
			days:   1,
			validate: func(result time.Time) bool {
				expected := time.Now().AddDate(0, 0, 1)
				return result.Day() == expected.Day()
			},
		},
		{
			name:   "subtract (negative values)",
			years:  -1,
			months: -1,
			days:   -1,
			validate: func(result time.Time) bool {
				expected := time.Now().AddDate(-1, -1, -1)
				return result.Year() == expected.Year() && result.Month() == expected.Month()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddDate(tt.years, tt.months, tt.days)
			if !tt.validate(result) {
				t.Errorf("AddDate(%d, %d, %d) returned invalid result", tt.years, tt.months, tt.days)
			}
		})
	}
}

func TestSecondsExpired(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int64
		expected bool
	}{
		{
			name:     "past time should be expired",
			seconds:  time.Now().Unix() - 100,
			expected: true,
		},
		{
			name:     "future time should not be expired",
			seconds:  time.Now().Unix() + 100,
			expected: false,
		},
		{
			name:     "zero should be expired",
			seconds:  0,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SecondsExpired(tt.seconds)
			if got != tt.expected {
				t.Errorf("SecondsExpired(%d) = %v, expected %v", tt.seconds, got, tt.expected)
			}
		})
	}
}

func TestMillisExpired(t *testing.T) {
	tests := []struct {
		name     string
		millis   int64
		expected bool
	}{
		{
			name:     "past time should be expired",
			millis:   time.Now().UnixMilli() - 1000,
			expected: true,
		},
		{
			name:     "future time should not be expired",
			millis:   time.Now().UnixMilli() + 1000,
			expected: false,
		},
		{
			name:     "zero should be expired",
			millis:   0,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MillisExpired(tt.millis)
			if got != tt.expected {
				t.Errorf("MillisExpired(%d) = %v, expected %v", tt.millis, got, tt.expected)
			}
		})
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		name           string
		timeA          time.Time
		timeB          time.Time
		expectedYear   int
		expectedMonth  int
		expectedDay    int
		expectedHour   int
		expectedMin    int
		expectedSec    int
	}{
		{
			name:          "same time",
			timeA:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			timeB:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedYear:  0,
			expectedMonth: 0,
			expectedDay:   0,
			expectedHour:  0,
			expectedMin:   0,
			expectedSec:   0,
		},
		{
			name:          "1 day difference",
			timeA:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			timeB:         time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			expectedYear:  0,
			expectedMonth: 0,
			expectedDay:   1,
			expectedHour:  0,
			expectedMin:   0,
			expectedSec:   0,
		},
		{
			name:          "1 month difference",
			timeA:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			timeB:         time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
			expectedYear:  0,
			expectedMonth: 1,
			expectedDay:   0,
			expectedHour:  0,
			expectedMin:   0,
			expectedSec:   0,
		},
		{
			name:          "1 year difference",
			timeA:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			timeB:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedYear:  1,
			expectedMonth: 0,
			expectedDay:   0,
			expectedHour:  0,
			expectedMin:   0,
			expectedSec:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, month, day, hour, min, sec := Diff(tt.timeA, tt.timeB)
			if year != tt.expectedYear || month != tt.expectedMonth || day != tt.expectedDay ||
				hour != tt.expectedHour || min != tt.expectedMin || sec != tt.expectedSec {
				t.Errorf("Diff() = (%d, %d, %d, %d, %d, %d), expected (%d, %d, %d, %d, %d, %d)",
					year, month, day, hour, min, sec,
					tt.expectedYear, tt.expectedMonth, tt.expectedDay,
					tt.expectedHour, tt.expectedMin, tt.expectedSec)
			}
		})
	}
}

func TestDiffExcluded(t *testing.T) {
	tests := []struct {
		name          string
		timeA         time.Time
		timeB         time.Time
		expectedYear  int
		expectedMonth int
		expectedDay   int
	}{
		{
			name:          "same day",
			timeA:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			timeB:         time.Date(2020, 1, 1, 23, 59, 59, 0, time.UTC),
			expectedYear:  0,
			expectedMonth: 0,
			expectedDay:   0,
		},
		{
			name:          "1 day difference",
			timeA:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			timeB:         time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			expectedYear:  0,
			expectedMonth: 0,
			expectedDay:   1,
		},
		{
			name:          "1 year difference",
			timeA:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			timeB:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedYear:  1,
			expectedMonth: 0,
			expectedDay:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, month, day := DiffExcluded(tt.timeA, tt.timeB)
			if year != tt.expectedYear || month != tt.expectedMonth || day != tt.expectedDay {
				t.Errorf("DiffExcluded() = (%d, %d, %d), expected (%d, %d, %d)",
					year, month, day,
					tt.expectedYear, tt.expectedMonth, tt.expectedDay)
			}
		})
	}
}

func TestTrimUnixTime(t *testing.T) {
	unixTime := int64(1577836800) // 2020-01-01 00:00:00 UTC
	result := TrimUnixTime(unixTime)

	if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
		t.Errorf("TrimUnixTime() should set hour, minute, second to 0, got %02d:%02d:%02d",
			result.Hour(), result.Minute(), result.Second())
	}
}

func TestAddSecMinHourUnixTime(t *testing.T) {
	unixTime := int64(1577836800) // 2020-01-01 00:00:00 UTC
	result := AddSecMinHourUnixTime(unixTime)

	if result.Hour() != 23 || result.Minute() != 59 || result.Second() != 59 {
		t.Errorf("AddSecMinHourUnixTime() should set hour:min:sec to 23:59:59, got %02d:%02d:%02d",
			result.Hour(), result.Minute(), result.Second())
	}
}

func TestAdjustStartEnd(t *testing.T) {
	startUnix := int64(1577836800) // 2020-01-01 00:00:00 UTC
	endUnix := int64(1577923199)   // 2020-01-01 23:59:59 UTC

	startResult, endResult := AdjustStartEnd(startUnix, endUnix)

	startTime := time.Unix(startResult, 0).UTC()
	if startTime.Hour() != 0 || startTime.Minute() != 0 || startTime.Second() != 0 {
		t.Errorf("AdjustStartEnd() start should be 00:00:00, got %02d:%02d:%02d",
			startTime.Hour(), startTime.Minute(), startTime.Second())
	}

	endTime := time.Unix(endResult, 0).UTC()
	if endTime.Hour() != 23 || endTime.Minute() != 59 || endTime.Second() != 59 {
		t.Errorf("AdjustStartEnd() end should be 23:59:59, got %02d:%02d:%02d",
			endTime.Hour(), endTime.Minute(), endTime.Second())
	}
}
