package utils

import (
	"os"
	"path"
	"time"
)

var (
	logsDir = "logs"
)

// Log2file write log to file
func Log2file(content, logName string) {
	var err error

	if logName == "" || len(logName) <= 0 {
		logName = "webhook-go.log"
	}
	dir, file := path.Split(logName)
	logsDir = path.Join(logsDir, dir)
	if _, err := os.Stat(logsDir); err != nil {
		err = os.MkdirAll(logsDir, 0711)
		if err != nil {
			return
		}
	}
	logPath := path.Join(logsDir, file)
	if _, err := os.Stat(logPath); err != nil {
		_, err = os.Create(logPath)
		if err != nil {
			return
		}
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	timeString := time.Now().Format("2006-01-02 15:04:05")
	_, err = f.WriteString("[" + timeString + "]" + "" + content)
	if err != nil {
		return
	}
	_, err = f.WriteString("\n")
	if err != nil {
		return
	}
}
