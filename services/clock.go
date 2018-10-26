package services

import (
	"time"
)

type Clock interface {
	GetTime() time.Time
}

type RealClock struct{}

func (c *RealClock) GetTime() time.Time {
	return time.Now()
}

type MockClock struct{}
