package internal

import (
	"github.com/syols/go-devops/config"
	"gopkg.in/yaml.v3"

	"testing"
)

func TestSettings(t *testing.T) {
	var configData = `
server:
  address:
    host: 0.0.0.0
    port: 80
  key: null

agent:
  poll_interval: 1s
  report_interval: 1s
  client_timeout: 1s

store:
  database: ""
  store_file: ""
  restore: true
  store_interval: 1s
`
	settings := config.Config{}
	if err := yaml.Unmarshal([]byte(configData), &settings); err != nil {
		t.Errorf(err.Error())
	}
}
