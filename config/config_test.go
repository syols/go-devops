package config

import (
	"testing"

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
	config := Config{}
	assert.NoError(t, yaml.Unmarshal([]byte(configData), &config))
}
