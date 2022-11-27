package main

import (
	"log"
	"os"

	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/app"
)

// @Title Agent API
// @Description Сервис сбора метрик
// @Version 0.1.0

// @Contact.email some@mail.com

// @BasePath /
// @Host 0.0.0.0:8080

//Example: go run -ldflags "-X main.buildVersion=%%" main.go
var buildVersion string
var buildDate string
var buildCommit string

func main() {
	log.SetOutput(os.Stdout)
	log.Printf("Build version: %s", config.ReplaceNoneValue(buildVersion))
	log.Printf("Build date: %s", config.ReplaceNoneValue(buildDate))
	log.Printf("Build commit: %s", config.ReplaceNoneValue(buildCommit))

	if server, err := app.NewServer(config.NewConfig()); err == nil {
		server.Run()
	}
}
