package eventserv

import (
	"bufio"
	"log"
	"os"
	"path/filepath"

	"github.com/CakeForKit/SysProtBtlnComp.git/internal/cnfg"
	"github.com/CakeForKit/SysProtBtlnComp.git/internal/utils"
)

type EventService interface {
	CreateReport(shortPathFile string)
}

func NewEventService(logger *log.Logger, raceCnfg *cnfg.RaceConfig) (EventService, error) {
	return &eventService{
		logger:   logger,
		raceCnfg: raceCnfg,
	}, nil
}

type eventService struct {
	logger   *log.Logger
	raceCnfg *cnfg.RaceConfig
}

func (s *eventService) CreateReport(shortPathFile string) {
	projectRoot := utils.GetProjectRoot()
	pathFile := filepath.Join(projectRoot, shortPathFile)
	file, err := os.Open(pathFile)
	if err != nil {
		s.logger.Fatalf("Ошибка открытия файла %s: %v", pathFile, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		s.logger.Print(line)
	}
	if err := scanner.Err(); err != nil {
		s.logger.Fatalf("Ошибка сканирования файла: %v", err)
	}
}
