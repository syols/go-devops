package settings

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Agent  AgentConfig  `yaml:"agent"`
	Store  StoreConfig  `yaml:"store"`
}

type ServerConfig struct {
	Address Address `yaml:"address"`
	Key     *string `yaml:"key,omitempty"`
}

type AgentConfig struct {
	PollInterval   time.Duration `yaml:"poll_interval"`
	ReportInterval time.Duration `yaml:"report_interval"`
	ClientTimeout  time.Duration `yaml:"client_timeout"`
}

type StoreConfig struct {
	DatabaseConnectionString *string       `yaml:"database,omitempty"`
	StoreFile                *string       `yaml:"store_file,omitempty"`
	Restore                  bool          `yaml:"restore"`
	StoreInterval            time.Duration `yaml:"store_interval"`
}

type Address struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

func NewConfig() (settings Config) {
	settings.setDefault("develop.yaml")
	settings.setFromOptions(newVariables().getOptions()...)
	return settings
}

func (s *Config) setFromOptions(options ...Option) {
	for _, fn := range options {
		fn(s)
	}
}

func (s *Config) setDefault(configPath string) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Print(err.Error())
	}

	if err := yaml.Unmarshal(file, s); err != nil {
		log.Print(err.Error())
	}
}

func (s *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Server.Address.Host, s.Server.Address.Port)
}

func (s *Config) String() (result string) {
	if marshal, err := yaml.Marshal(s); err != nil {
		result = string(marshal)
	}
	return
}
