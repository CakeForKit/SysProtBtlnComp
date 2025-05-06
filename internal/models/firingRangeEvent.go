package models

import "time"

type firingRangeEvent struct {
	timestamp time.Time
	id        int
	num       int
}

func NewFiringRangeEvent(timestamp time.Time, id int, num int) (Event, error) {
	return &firingRangeEvent{
		timestamp: timestamp,
		id:        id,
		num:       num,
	}, nil
}

func (e *firingRangeEvent) GetTimestamp() time.Time {
	return e.timestamp
}

func (e *firingRangeEvent) GetID() int {
	return e.id
}

func (e *firingRangeEvent) GetStartTime() time.Time {
	panic("it's not startEvent!")
}

func (e *firingRangeEvent) GetNumber() int {
	return e.num
}

func (e *firingRangeEvent) GetComment() string {
	panic("it's not commentEvent!")
}
