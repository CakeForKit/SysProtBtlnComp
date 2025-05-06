package eventserv

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/CakeForKit/SysProtBtlnComp.git/internal/cnfg"
	"github.com/CakeForKit/SysProtBtlnComp.git/internal/models"
)

var (
	ErrFormat = errors.New("error format")
)

type EventService interface {
	CreateReport(shortPathFile string, outputPathFile string)
}

func NewEventService(logger *log.Logger, raceCnfg *cnfg.RaceConfig) (EventService, error) {
	return &eventService{
		logger:   logger,
		raceCnfg: raceCnfg,
	}, nil
}

type eventService struct {
	logger      *log.Logger
	raceCnfg    *cnfg.RaceConfig
	competitors map[int]models.Competitor
}

func (s *eventService) CreateReport(inputPathFile string, outputPathFile string) {
	file, err := os.Open(inputPathFile)
	if err != nil {
		log.Fatalf("Ошибка открытия файла %s: %v", inputPathFile, err)
	}
	defer file.Close()

	outfile, err := os.OpenFile(outputPathFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Ошибка открытия файла %s: %v", inputPathFile, err)
	}
	defer outfile.Close()

	scanner := bufio.NewScanner(file)
	s.competitors = make(map[int]models.Competitor, 0)
	for scanner.Scan() {
		line := scanner.Text()
		event, competitorID, err := s.parseLine(line)
		if err != nil {
			log.Fatalf("CreateReport: %v", err)
		}
		if val, ok := s.competitors[competitorID]; ok {
			err = val.AddEvent(event, s.raceCnfg, s.logger)
			if err != nil {
				log.Fatalf("CreateReport: %v", err)
			}
		} else {
			c, err := models.NewCompetitor(competitorID)
			if err != nil {
				log.Fatalf("NewCompetitor: %v", err)
			}
			c.AddEvent(event, s.raceCnfg, s.logger)
			s.competitors[competitorID] = c
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Ошибка сканирования файла: %v", err)
	}

	result := make([]models.Competitor, 0, len(s.competitors))
	for _, competitor := range s.competitors {
		result = append(result, competitor)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].GetTotalTime() < result[j].GetTotalTime()
	})

	for _, value := range result {
		stat, err := value.GetStatistic(s.raceCnfg)
		if err != nil {
			log.Fatalf("GetStatisticText: %v", err)
		}
		// fmt.Printf("%s\n", stat)
		outfile.WriteString(fmt.Sprintf("%s\n", stat))
	}

}

func (s *eventService) parseLine(line string) (models.Event, int, error) {
	parts := strings.SplitN(line, "]", 2)
	if len(parts) < 2 {
		return nil, 0, fmt.Errorf("cannot find []: %v", ErrFormat)
	}

	timestr := strings.TrimSpace(parts[0][1:])
	timestamp, err := time.Parse("15:04:05.000", timestr)
	if err != nil {
		return nil, 0, fmt.Errorf("err parse time: %v", err)
	}

	restParts := strings.Fields(parts[1])
	if len(parts) < 2 {
		return nil, 0, fmt.Errorf("params must be min 2: %v", ErrFormat)
	}

	eventID, err := strconv.Atoi(restParts[0])
	if err != nil {
		return nil, 0, fmt.Errorf("error read eventID: %w", err)
	}
	competitorID, err := strconv.Atoi(restParts[1])
	if err != nil {
		return nil, 0, fmt.Errorf("error read competitorID: %w", err)
	}

	var event models.Event
	if len(restParts) == 2 {
		event, err = models.NewSimpleEvent(timestamp, eventID)
		if err != nil {
			return nil, 0, fmt.Errorf("err create simple event: %v", err)
		}
	} else {
		if eventID == 2 {
			startTime, err := time.Parse("15:04:05.000", restParts[2])
			if err != nil {
				return nil, 0, fmt.Errorf("err parse start time: %v", err)
			}
			event, err = models.NewStartEvent(timestamp, eventID, startTime)
			if err != nil {
				return nil, 0, fmt.Errorf("err create start event: %v", err)
			}
		} else if eventID == 5 || eventID == 6 {
			num, err := strconv.Atoi(restParts[2])
			if err != nil {
				return nil, 0, fmt.Errorf("err parse num: %v", err)
			}
			if eventID == 5 {
				event, err = models.NewFiringRangeEvent(timestamp, eventID, num)
			} else {
				event, err = models.NewTargetEvent(timestamp, eventID, num)
			}
			if err != nil {
				return nil, 0, fmt.Errorf("err create number event: %v", err)
			}
		} else if eventID == 11 {
			comment := strings.Join(restParts[2:], " ")
			event, err = models.NewCommentEvent(timestamp, eventID, comment)
			if err != nil {
				return nil, 0, fmt.Errorf("err create comment event: %v", err)
			}
		} else {
			return nil, 0, fmt.Errorf("err extra param")
		}
	}

	return event, competitorID, nil
}
