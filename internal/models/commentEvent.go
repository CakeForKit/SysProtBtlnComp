package models

import "time"

type commentEvent struct {
	timestamp time.Time
	id        int
	comment   string
}

func NewCommentEvent(timestamp time.Time, id int, comment string) (Event, error) {
	return &commentEvent{
		timestamp: timestamp,
		id:        id,
		comment:   comment,
	}, nil
}

func (e *commentEvent) GetTimestamp() time.Time {
	return e.timestamp
}

func (e *commentEvent) GetID() int {
	return e.id
}

func (e *commentEvent) GetStartTime() time.Time {
	panic("it's not startEvent!")
}

func (e *commentEvent) GetNumber() int {
	panic("it's not numberEvent!")
}

func (e *commentEvent) GetComment() string {
	return e.comment
}
