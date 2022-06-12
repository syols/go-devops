package main

import (
	"github.com/syols/go-devops/internal"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	server := internal.NewServer(internal.NewSettings())
	server.Run()
}
