package main

import (
	"github.com/syols/go-devops/internal"
	"github.com/syols/go-devops/internal/settings"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	sets := settings.NewSettings()
	server := internal.NewServer(sets)
	server.Run()
}
