package utils

import "time"

const TimeLayout = "2006-01-02T15:04:05.000Z"

func Time2Str(t time.Time) string {
	return t.Format(TimeLayout)
}

func NowStr() string {
	return Time2Str(time.Now())
}

func FromNow(delay time.Duration) time.Time {
	return time.Now().Add(delay)
}
