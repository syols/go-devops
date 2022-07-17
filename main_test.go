package main

import (
	"github.com/stretchr/testify/require"
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/app"
	"net"
	"testing"
	"time"
)

func settings(t *testing.T) (config.Config, error) {
	list, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	port := list.Addr().(*net.TCPAddr).Port
	err = list.Close()
	require.NoError(t, err)
	settings := config.Config{
		Server: config.ServerConfig{
			Address: config.Address{
				Host: "0.0.0.0",
				Port: uint16(port),
			},
		},
		Store: config.StoreConfig{
			StoreInterval: time.Second * 10,
		},
	}
	return settings, err
}

func TestStartServer(t *testing.T) {
	sets, err := settings(t)
	require.NoError(t, err)
	server, err := app.NewServer(sets)
	go server.Run()
}
