package utils

import (
	"path/filepath"
	"runtime"
)

func GetProjectRoot() string {
	_, currentFile, _, _ := runtime.Caller(0) // Получаем путь к текущему файлу
	projectRoot := filepath.Join(filepath.Dir(currentFile), "..", "..")
	return projectRoot
}
