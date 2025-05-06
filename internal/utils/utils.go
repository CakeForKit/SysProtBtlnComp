package utils

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/CakeForKit/SysProtBtlnComp.git/internal/models"
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

func CreateComment(competitorID int, event models.Event) string {
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

func GetProjectRoot() string {
	_, currentFile, _, _ := runtime.Caller(0) // Получаем путь к текущему файлу
	projectRoot := filepath.Join(filepath.Dir(currentFile), "..", "..")
	return projectRoot
}
