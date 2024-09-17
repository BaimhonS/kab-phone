package utils

import "time"

func ParseToLocalTime(timeUTC time.Time) time.Time {
	return timeUTC.In(time.Local)
}

func GetStartOfDay() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, now.Location())
}

func GetEndOfDay() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, now.Location())
}
