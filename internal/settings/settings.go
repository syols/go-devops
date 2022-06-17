package settings

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

type Settings struct {
	Server ServerSettings `yaml:"server"`
	Agent  AgentSettings  `yaml:"agent"`
	Store  StoreSettings  `yaml:"store"`
}

type ServerSettings struct {
	Address Address `yaml:"address"`
	Key     *string `yaml:"key,omitempty"`
}

type AgentSettings struct {
	PollInterval   time.Duration `yaml:"poll_interval"`
	ReportInterval time.Duration `yaml:"report_interval"`
	ClientTimeout  time.Duration `yaml:"client_timeout"`
}

type StoreSettings struct {
	DatabaseConnectionString *string       `yaml:"database,omitempty"`
	StoreFile                *string       `yaml:"store_file,omitempty"`
	Restore                  bool          `yaml:"restore"`
	StoreInterval            time.Duration `yaml:"store_interval"`
}

type Address struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

func NewSettings() (settings Settings) {
	log.Print("Settings:")
	settings.setDefault("configs/default.yml")
	settings.setFromOptions(newVariables().getOptions()...)
	log.Printf(settings.String())
	return settings
}

func (s *Settings) setFromOptions(options ...Option) {
	for _, fn := range options {
		fn(s)
	}
}

func (s *Settings) setDefault(configPath string) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("Read config error")
	}

	if err := yaml.Unmarshal(file, s); err != nil {
		log.Printf("Config load error")
	}
}

func (s *Settings) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Server.Address.Host, s.Server.Address.Port)
}

func (s *Settings) String() (result string) {
	if marshal, err := yaml.Marshal(s); err != nil {
		result = string(marshal)
	}
	return
}
