package main

import (
	"github.com/syols/go-devops/internal"
	"github.com/syols/go-devops/internal/settings"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Printf("Server")

	server := internal.NewServer(settings.NewSettings())
	server.Run()
}
