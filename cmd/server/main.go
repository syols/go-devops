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

func main() {
	log.SetOutput(os.Stdout)
	if server, err := app.NewServer(config.NewConfig()); err == nil {
		server.Run()
	}
}
