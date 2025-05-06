package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/CakeForKit/SysProtBtlnComp.git/internal/cnfg"
	"github.com/CakeForKit/SysProtBtlnComp.git/internal/eventserv"
	"github.com/CakeForKit/SysProtBtlnComp.git/internal/utils"
)

func main() {
	// Logger
	projectRoot := utils.GetProjectRoot()
	logFilePath := filepath.Join(projectRoot, "log/log.log")
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Ошибка открытия файла логов:", err)
	}
	defer logFile.Close()
	// ----------

	// Config
	raceConfig, err := cnfg.LoadRaceConfig("./configs/")
	if err != nil {
		fmt.Printf("ERR %v\n", err)
	}
	fmt.Printf("%+v\n", raceConfig)
	// -----------

	logger := log.New(logFile, "", 0)
	// timestamp := time.Now().Format("15:04:05.000")
	// id := 1
	// logger.Printf("[%s] The competitor(%d) registered", timestamp, id)
	raceServ, err := eventserv.NewEventService(logger, raceConfig)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	raceServ.CreateReport("input_data/test")
}
