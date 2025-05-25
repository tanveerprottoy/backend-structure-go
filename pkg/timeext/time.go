package timeext

import "time"

// NowUnixMilli returns the current time in Unix milliseconds
func NowUnixMilli() int64 {
	return time.Now().UnixMilli()
}

// AddDate returns the time corresponding to adding the given number of years, months, and days to t.
func AddDate(years int, months int, days int) time.Time {
	return time.Now().AddDate(years, months, days)
}

// SecondsExpired returns true if the given seconds have expired
func SecondsExpired(seconds int64) bool {
	return time.Now().Unix() > seconds
}

// MillisExpired returns true if the given milliseconds have expired
func MillisExpired(millis int64) bool {
	return time.Now().UnixMilli() > millis
}

// Diff returns the difference between two times
func Diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}
	return
}

// DiffExcluded returns the difference between two times, excluding the seconds, minutes, and hours
func DiffExcluded(a, b time.Time) (year, month, day int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)

	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}
	return
}

func TrimUnixTime(u int64) time.Time {
	t := time.Unix(u, 0).UTC()
	// remove seconds, minutes, and hours
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func AddSecMinHourUnixTime(u int64) time.Time {
	t := time.Unix(u, 0).UTC()
	// remove seconds, minutes, and hours
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

func AdjustStartEnd(start, end int64) (int64, int64) {
	t0 := TrimUnixTime(start)
	t1 := AddSecMinHourUnixTime(end)
	return t0.Unix(), t1.Unix()
}
