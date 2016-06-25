package clock

import "time"

type Clock interface {
	Now() time.Time
}

func NewClock() Clock {
	return realClock{}
}

type realClock struct{}

func (c realClock) Now() time.Time {
	return time.Now()
}
