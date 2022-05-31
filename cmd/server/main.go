package main

import (
	"github.com/syols/go-devops/internal"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	var settings internal.Settings
	settings.LoadSettings(internal.ConfigPath)
	newServer := internal.NewServer(settings)
	newServer.Run()
}
