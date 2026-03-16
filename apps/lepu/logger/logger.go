package logger

import (
	"io"
	"log"
	"os"
)

var Logger *log.Logger
var LogWriter io.Writer

func InitLogger() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	LogWriter = io.MultiWriter(os.Stdout, file)
	Logger = log.New(LogWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}