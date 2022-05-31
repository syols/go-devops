package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

const ConfigPath = "configs/default.yml"
const ConfigLoadErrorMessage = "Config load error"

type Settings struct {
	Address Address `yaml:"address"`
	Agent   Agent   `yaml:"agent"`
	Metrics Metrics `yaml:"metrics"`
}

type Address struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Agent struct {
	PollInterval   time.Duration `yaml:"poll_interval"`
	ReportInterval time.Duration `yaml:"report_interval"`
	ClientTimeout  time.Duration `yaml:"client_timeout"`
}

type CustomMetric struct {
	Name      string `yaml:"name"`
	TypeAlias string `yaml:"type_alias"`
}

type Metrics struct {
	RuntimeMetrics []string `yaml:"runtime"`
}

func (settings *Settings) LoadSettings(configPath string) {
	if file, err := ioutil.ReadFile(ConfigPath); err == nil {
		if err = yaml.Unmarshal(file, settings); err == nil {
			return
		}
	}
	log.Fatalf(ConfigLoadErrorMessage)
}

func (settings *Settings) GetAddress() string {
	return fmt.Sprintf("%s:%d", settings.Address.Host, settings.Address.Port)
}
