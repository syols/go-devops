package main

import (
	"github.com/syols/go-devops/internal/server"
	"github.com/syols/go-devops/internal/utils"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	var settings utils.Settings
	settings.LoadSettings()
	newServer := server.NewServer(settings)
	newServer.Run()
}
