package models

import (
	"time"
)

type Event interface {
	GetTimestamp() time.Time
	GetID() int
	GetStartTime() time.Time
	GetNumber() int
	GetComment() string
}
