package internal

import (
	"github.com/syols/go-devops/internal/settings"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestSettings(t *testing.T) {
	var config = `
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
	sets := settings.Config{}
	if err := yaml.Unmarshal([]byte(config), &sets); err != nil {
		t.Errorf(err.Error())
	}
}
