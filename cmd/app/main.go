package main

import (
	"log"
	"os"

	"github.com/CakeForKit/SysProtBtlnComp.git/internal/cnfg"
	"github.com/CakeForKit/SysProtBtlnComp.git/internal/eventserv"
)

func main() {
	// Logger
	logFilePath := "log/log.log"
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Ошибка открытия файла логов:", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", 0)
	// logger := log.New(os.Stdout, "", 0)
	// ----------

	// Config
	raceConfig, err := cnfg.LoadRaceConfig("./configs/")
	if err != nil {
		log.Fatalf("ERR %v\n", err)
	}
	// -----------

	raceServ, err := eventserv.NewEventService(logger, raceConfig)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	raceServ.CreateReport("input_data/events", "log/result")

}
