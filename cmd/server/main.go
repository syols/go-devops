package main

import (
	"log"
	"os"

	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/app"
)

//Example go run -ldflags "-X main.buildVersion=%%" main.go
var buildVersion string
var buildDate string
var buildCommit string

// @Title Agent API
// @Description Сервис сбора метрик
// @Version 0.1.0

// @Contact.email some@mail.com

// @BasePath /
// @Host 0.0.0.0:8080

func main() {
	log.SetOutput(os.Stdout)

	checkNone := func(value string) string {
		if len(value) == 0 {
			return "N/A"
		}
		return value
	}

	log.Printf("Build version: %s", checkNone(buildVersion))
	log.Printf("Build date: %s", checkNone(buildDate))
	log.Printf("Build commit: %s", checkNone(buildCommit))

	if server, err := app.NewServer(config.NewConfig()); err == nil {
		server.Run()
	}
}
