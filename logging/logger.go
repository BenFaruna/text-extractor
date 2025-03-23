package logging

import (
	"io"
	"log"
	"log/slog"
	"os"
)

var logger *slog.Logger
var logFile io.WriteCloser

func init() {
	var filename string
	fileName := os.Getenv("LOG_FILE")
	if fileName == "" {
		filename = "app.log"
	} else {
		filename = fileName
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logFile = file

	textHandler := slog.NewTextHandler(file, nil)
	logger = slog.New(textHandler)
	slog.SetDefault(logger)
}

func GetLogger() *slog.Logger {
	return logger
}

func CloseLogger() {
	err := logFile.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
}
