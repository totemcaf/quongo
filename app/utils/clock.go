package utils

import (
	"log"
	"time"
)

type Clock interface {
	Now() time.Time
}

type productionClock struct {
}

func ProductionClock() Clock {
	return &productionClock{}
}

func (c *productionClock) Now() time.Time {
	return time.Now()
}

type fixedClock struct {
	fixedTime time.Time
}

func FixedClockAtNow() Clock {
	return FixedClock(time.Now())
}

// FixedClockAt at time parse from YYYY-MM-DD HH:MM:SS
// Mon Jan 2 15:04:05 -0700 MST 2006
func FixedClockAt(fixedTimeStr string) Clock {
	fixedTime, err := time.Parse("2006-01-02 15:04:05", fixedTimeStr)
	if err != nil {
		log.Fatalf("Invalid time string: '{}'", fixedTimeStr)
	}
	return FixedClock(fixedTime)
}

func FixedClock(fixedTime time.Time) Clock {
	return &fixedClock{fixedTime}
}

func (c *fixedClock) Now() time.Time {
	return c.fixedTime
}
