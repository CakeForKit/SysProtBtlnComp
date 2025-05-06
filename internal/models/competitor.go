package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/CakeForKit/SysProtBtlnComp.git/internal/cnfg"
)

var (
	eventComments = map[int]string{
		1:  "[%s] The competitor(%d) registered",
		2:  "[%s] The start time for the competitor(%d) was set by a draw to %s",
		3:  "[%s] The competitor(%d) is on the start line",
		4:  "[%s] The competitor(%d) has started",
		5:  "[%s] The competitor(%d) is on the firing range(%d)",
		6:  "[%s] The target(%d) has been hit by competitor(%d)",
		7:  "[%s] The competitor(%d) left the firing range",
		8:  "[%s] The competitor(%d) entered the penalty laps",
		9:  "[%s] The competitor(%d) left the penalty laps",
		10: "[%s] The competitor(%d) ended the main lap",
		11: "[%s] The competitor(%d) can't continue: %s",
		32: "[%s] The competitor(%d) is disqualified",
		33: "[%s] The competitor(%d) has finished",
	}
)

type Competitor interface {
	AddEvent(event Event, raceCnfg *cnfg.RaceConfig, logger *log.Logger) error
	GetEvents() []Event
	GetStatistic(raceCnfg *cnfg.RaceConfig) (string, error)
	GetTotalTime() time.Duration
}

func NewCompetitor(id int) (Competitor, error) {
	return &competitor{
		id:                   id,
		events:               make([]Event, 0),
		onPenaltyLap:         false,
		registered:           false,
		setTimeStart:         false,
		scheduledTimeStart:   time.Time{},
		onStartLine:          false,
		isStarted:            false,
		actualTimeStart:      time.Time{},
		status:               "",
		hits:                 0,
		startTimePenaltyLap:  time.Time{},
		cntPenaltyLaps:       0,
		totalTimePenaltyLaps: 0,
		startTimeCurLap:      time.Time{},
		timeEachLap:          make([]time.Duration, 0),
	}, nil
}

type competitor struct {
	id                   int
	events               []Event
	onPenaltyLap         bool
	registered           bool
	setTimeStart         bool
	scheduledTimeStart   time.Time
	onStartLine          bool
	isStarted            bool
	actualTimeStart      time.Time
	status               string
	hits                 int
	startTimePenaltyLap  time.Time
	cntPenaltyLaps       int
	totalTimePenaltyLaps time.Duration
	startTimeCurLap      time.Time
	timeEachLap          []time.Duration
}

func (c *competitor) AddEvent(event Event, raceCnfg *cnfg.RaceConfig, logger *log.Logger) error {
	if c.status != "" {
		return fmt.Errorf("already out of competition")
	}
	elen := len(c.events)
	if elen > 0 && c.events[elen-1].GetTimestamp().After(event.GetTimestamp()) {
		return fmt.Errorf("AddEvent: new event %d before old event %d", event.GetID(), c.events[elen-1].GetID())
	}
	id := event.GetID()
	switch id {
	case 1: // The competitor registered
		if c.registered {
			return fmt.Errorf("already registered")
		}
		c.registered = true
	case 2: //The start time was set by a draw
		if c.setTimeStart {
			return fmt.Errorf("already setTimeStart")
		}
		c.setTimeStart = true
		c.scheduledTimeStart = event.GetStartTime()
	case 3: // The competitor is on the start line
		if c.onStartLine {
			return fmt.Errorf("already onStartLine")
		}
		c.onStartLine = true
	case 4: // The competitor has started
		if c.isStarted {
			return fmt.Errorf("already isStarted")
		}
		c.isStarted = true
		c.actualTimeStart = event.GetTimestamp()
		if c.actualTimeStart.Before(c.scheduledTimeStart) ||
			c.actualTimeStart.After(c.scheduledTimeStart.Add(raceCnfg.StartDelta)) {
			c.status = "NotStarted"
		}
		c.startTimeCurLap = c.actualTimeStart
	case 5: // The competitor is on the firing range
	case 6: // The target has been hit
		c.hits++
	case 7: // The competitor left the firing range
	case 8: // The competitor entered the penalty laps
		if c.onPenaltyLap {
			return fmt.Errorf("AddEvent: already on Penalty Lap")
		}
		c.onPenaltyLap = true
		c.startTimePenaltyLap = event.GetTimestamp()
		c.cntPenaltyLaps++
	case 9: // The competitor left the penalty laps
		if !c.onPenaltyLap {
			return fmt.Errorf("AddEvent: wasnot on Penalty Lap")
		}
		c.onPenaltyLap = false
		c.totalTimePenaltyLaps += event.GetTimestamp().Sub(c.startTimePenaltyLap)
		// c.startTimeCurLap = event.GetTimestamp()
	case 10: // The competitor ended the main lap
		timeLap := event.GetTimestamp().Sub(c.startTimeCurLap)
		c.timeEachLap = append(c.timeEachLap, timeLap)
		c.startTimeCurLap = event.GetTimestamp()
		if len(c.timeEachLap) == raceCnfg.Laps {
			newEvent, err := NewSimpleEvent(event.GetTimestamp(), 33)
			if err != nil {
				return fmt.Errorf("error NewSimpleEvent")
			}
			logger.Printf("%s\n", c.createComment(c.id, event))
			logger.Printf("%s\n", c.createComment(c.id, newEvent))
			return nil
		}
	case 11: // The competitor can`t continue
		c.status = "NotFinished"
		newEvent, err := NewSimpleEvent(event.GetTimestamp(), 32)
		if err != nil {
			return fmt.Errorf("error NewSimpleEvent")
		}
		c.events = append(c.events, event)
		c.events = append(c.events, newEvent)
		logger.Printf("%s\n", c.createComment(c.id, event))
		logger.Printf("%s\n", c.createComment(c.id, newEvent))
		return nil
	}
	c.events = append(c.events, event)
	logger.Printf("%s\n", c.createComment(c.id, event))
	return nil
}

func (c *competitor) GetEvents() []Event {
	return c.events
}

func (c *competitor) GetTotalTime() time.Duration {
	return c.totalTimePenaltyLaps
}

func formatDuration(d time.Duration) string {
	t := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	t = t.Add(d)
	return t.Format("15:04:05.000")
}

func (c *competitor) GetStatistic(raceCnfg *cnfg.RaceConfig) (string, error) {
	if c.status == "" {
		if !c.registered || !c.setTimeStart || !c.onStartLine || !c.isStarted {
			c.status = "NotStarted"
		} else {
			c.status = formatDuration(c.totalTimePenaltyLaps)
		}
	}
	statLaps := make([]string, 0)
	for _, timeLap := range c.timeEachLap {
		avgSpeed := float64(raceCnfg.LapLen) / timeLap.Seconds()
		statLaps = append(statLaps, fmt.Sprintf("{%s, %.3f}", formatDuration(timeLap), avgSpeed))
	}

	shots := 5 * raceCnfg.FiringLines
	var avgSpeedPenaltyLaps float64
	if c.cntPenaltyLaps == 0 {
		avgSpeedPenaltyLaps = 0
	} else {
		avgSpeedPenaltyLaps = float64(c.cntPenaltyLaps*raceCnfg.PenaltyLen) / c.totalTimePenaltyLaps.Seconds()
	}

	statText := fmt.Sprintf(
		"[%s] %d %s {%s, %.3f} %d/%d",
		c.status,
		c.id,
		"["+strings.Join(statLaps, ", ")+"]",
		formatDuration(c.totalTimePenaltyLaps),
		avgSpeedPenaltyLaps,
		c.hits, shots)
	return statText, nil
}

func (c *competitor) createComment(competitorID int, event Event) string {
	id := event.GetID()
	timestamp := event.GetTimestamp().Format("15:04:05.000")
	switch id {
	case 2:
		return fmt.Sprintf(eventComments[id], timestamp, competitorID, event.GetStartTime().Format("15:04:05.000"))
	case 5:
		return fmt.Sprintf(eventComments[id], timestamp, competitorID, event.GetNumber())
	case 6:
		return fmt.Sprintf(eventComments[id], timestamp, event.GetNumber(), competitorID)
	case 11:
		return fmt.Sprintf(eventComments[id], timestamp, competitorID, event.GetComment())
	default:
		return fmt.Sprintf(eventComments[id], timestamp, competitorID)
	}
}
