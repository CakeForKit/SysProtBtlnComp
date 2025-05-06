package models

import "time"

type startEvent struct {
	timestamp time.Time
	id        int
	startTime time.Time
}

func NewStartEvent(timestamp time.Time, id int, startTime time.Time) (Event, error) {
	return &startEvent{
		timestamp: timestamp,
		id:        id,
		startTime: startTime,
	}, nil
}

func (e *startEvent) GetTimestamp() time.Time {
	return e.timestamp
}

func (e *startEvent) GetID() int {
	return e.id
}

func (e *startEvent) GetStartTime() time.Time {
	return e.startTime
}

func (e *startEvent) GetNumber() int {
	panic("it's not numberEvent!")
}

func (e *startEvent) GetComment() string {
	panic("it's not commentEvent!")
}
