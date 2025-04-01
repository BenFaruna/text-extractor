package logging

import (
	"io"
	"log"
	"os"
)

var (
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
)

// var logger *slog.Logger
var logFile io.WriteCloser

func init() {
	var filename string
	fileName := os.Getenv("LOG_FILE")
	if fileName == "" {
		filename = "app.log"
	} else {
		filename = fileName
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logFile = file

	log.SetOutput(logFile)
	flags := log.Ldate | log.Ltime | log.Lshortfile

	WarnLogger = log.New(logFile, "WARN: ", flags)
	InfoLogger = log.New(logFile, "INFO: ", flags)
	ErrorLogger = log.New(logFile, "ERROR: ", flags)
}

func CloseLogger() {
	err := logFile.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
}
