package internal

import (
	"github.com/syols/go-devops/config"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestSettings(t *testing.T) {
	var configData = `
address:
  host: 0.0.0.0
  port: 8080

agent:
  poll_interval: 1000
  report_interval: 10000
  client_timeout: 3000

metrics:
  runtime:
    - Alloc
    - BuckHashSys
`
	settings := config.Config{}
	if err := yaml.Unmarshal([]byte(configData), &settings); err != nil {
		t.Errorf(err.Error())
	}
}
