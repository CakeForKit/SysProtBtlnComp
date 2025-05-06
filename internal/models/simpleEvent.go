package models

import "time"

type simpleEvent struct {
	timestamp time.Time
	id        int
}

func NewSimpleEvent(timestamp time.Time, id int) (Event, error) {
	return &simpleEvent{
		timestamp: timestamp,
		id:        id,
	}, nil
}

func (e *simpleEvent) GetTimestamp() time.Time {
	return e.timestamp
}

func (e *simpleEvent) GetID() int {
	return e.id
}

func (e *simpleEvent) GetStartTime() time.Time {
	panic("it's not startEvent!")
}

func (e *simpleEvent) GetNumber() int {
	panic("it's not numberEvent!")
}

func (e *simpleEvent) GetComment() string {
	panic("it's not commentEvent!")
}
