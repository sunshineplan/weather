package unit

import "time"

type UnixTime int64

func (t UnixTime) Time() time.Time {
	return time.Unix(int64(t), 0)
}

func (t UnixTime) Date() string {
	return t.Time().Format("2006-01-02")
}

func (t UnixTime) Until() time.Duration {
	return time.Until(t.Time().AddDate(0, 0, 1)).Truncate(24 * time.Hour)
}

func (t UnixTime) Before(date time.Time) bool {
	year, month, day := date.Date()
	y, m, d := t.Time().Date()
	if year == y {
		if month == m {
			return day > d
		}
		return month > m
	}
	return year > y
}

func (t UnixTime) IsExpired() bool {
	return t.Before(time.Now())
}

func (t UnixTime) Weekday() string {
	return t.Time().Weekday().String()[:3]
}
