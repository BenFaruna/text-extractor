package main

import (
	"github.com/BenFaruna/text-extractor/logging"
	_ "github.com/BenFaruna/text-extractor/logging"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Project initialized", "file", "main.go", "line", 10)

	logging.CloseLogger()
}
