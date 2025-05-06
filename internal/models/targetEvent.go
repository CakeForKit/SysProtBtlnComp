package models

import (
	"fmt"
	"time"
)

// firingRange or target
type targetEvent struct {
	timestamp time.Time
	id        int
	num       int
}

func NewTargetEvent(timestamp time.Time, id int, num int) (Event, error) {
	if num < 1 || num > 5 {
		return nil, fmt.Errorf("error target num")
	}
	return &targetEvent{
		timestamp: timestamp,
		id:        id,
		num:       num,
	}, nil
}

func (e *targetEvent) GetTimestamp() time.Time {
	return e.timestamp
}

func (e *targetEvent) GetID() int {
	return e.id
}

func (e *targetEvent) GetStartTime() time.Time {
	panic("it's not startEvent!")
}

func (e *targetEvent) GetNumber() int {
	return e.num
}

func (e *targetEvent) GetComment() string {
	panic("it's not commentEvent!")
}
