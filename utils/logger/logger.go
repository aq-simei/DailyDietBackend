package logger

import (
	"fmt"
	"time"
)

type LogType string

const (
	INFO    LogType = "INFO"
	ERROR   LogType = "ERROR"
	WARNING LogType = "WARNING"
	DEBUG   LogType = "DEBUG"
)

var colorMap = map[LogType]string{
	INFO:    "\033[32m", // Green
	ERROR:   "\033[31m", // Red
	WARNING: "\033[33m", // Yellow
	DEBUG:   "\033[36m", // Cyan
}

func Log(logType LogType, message string) {
	colorCode := colorMap[logType]
	resetColor := "\033[0m"
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	fmt.Printf("%s[%s] %s: %s%s\n",
		colorCode,
		timestamp,
		logType,
		message,
		resetColor,
	)
}
