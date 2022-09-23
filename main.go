package main

import (
	"int-service/app"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.SetFormatter(&logrus.JSONFormatter{})

	app.Initialize("2002", logger)
}
