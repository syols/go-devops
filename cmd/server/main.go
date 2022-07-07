package main

import (
	"github.com/syols/go-devops/internal"
	"github.com/syols/go-devops/internal/settings"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	server := internal.NewServer(settings.NewConfig())
	server.Run()
}
