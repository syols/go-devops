package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	var configData = `
server:
  address:
    host: 0.0.0.0
    port: 8080
  key: null

agent:
  poll_interval: 2s
  report_interval: 3s
  client_timeout: 1s

store:
  database: null #"postgres://postgres:postgres@localhost/postgres?sslmode=disable"
  store_file: "/tmp/devops-metrics-db.json"
  restore: true
  store_interval: 5m0s
`
	storeFile := "/tmp/devops-metrics-db.json"
	config := Config{}
	expected := Config{
		Server: ServerConfig{
			Address: Address{
				Host: "0.0.0.0",
				Port: uint16(8080),
			},
		},
		Store: StoreConfig{
			DatabaseConnectionString: nil,
			StoreFile:                &storeFile,
			Restore:                  true,
			StoreInterval:            5 * time.Minute,
		},
		Agent: AgentConfig{
			PollInterval:   2 * time.Second,
			ReportInterval: 3 * time.Second,
			ClientTimeout:  time.Second,
		},
	}

	assert.NoError(t, yaml.Unmarshal([]byte(configData), &config))
	assert.Equal(t, config, expected)
}

func TestDefaultConfig(t *testing.T) {
	config := Config{}
	err := config.setDefault("../develop.yml")
	assert.NoError(t, err)
}

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	assert.Equal(t, config.Agent.ReportInterval, time.Duration(0))
}

func TestStringConfig(t *testing.T) {
	config := NewConfig()
	assert.Equal(t, config.String(), "")
}
