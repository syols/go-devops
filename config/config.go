package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"time"
)

type Option func(s *Config)

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
	settings.setDefault("develop.yml")
	settings.setFromOptions(newVariables().Options()...)
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

func (s *Config) Address() string {
	return fmt.Sprintf("%s:%d", s.Server.Address.Host, s.Server.Address.Port)
}

func (s *Config) String() (result string) {
	if marshal, err := yaml.Marshal(s); err != nil {
		result = string(marshal)
	}
	return
}

func WithAddress(address string) Option {
	return func(s *Config) {
		if host, port, err := net.SplitHostPort(address); err == nil {
			if port, err := strconv.ParseUint(port, 0, 16); err == nil {
				s.Server.Address.Host = host
				s.Server.Address.Port = uint16(port)
			}
		}
	}
}

func WithRestore(value string) Option {
	return func(s *Config) {
		val, err := strconv.ParseBool(value)
		if err != nil {
			log.Print(err.Error())
		}
		s.Store.Restore = val
	}
}

func WithStoreFile(value string) Option {
	return func(s *Config) {
		s.Store.StoreFile = &value
	}
}

func WithPollInterval(value string) Option {
	return func(s *Config) {
		val, err := time.ParseDuration(value)
		if err != nil {
			log.Print(err.Error())
		}
		s.Agent.PollInterval = val
	}
}

func WithReportInterval(value string) Option {
	return func(s *Config) {
		val, err := time.ParseDuration(value)
		if err != nil {
			log.Print(err.Error())
		}
		s.Agent.ReportInterval = val
	}
}

func WithClientTimeout(value string) Option {
	return func(s *Config) {
		val, err := time.ParseDuration(value)
		if err != nil {
			log.Print(err.Error())
		}
		s.Agent.ClientTimeout = val
	}
}

func WithStoreInterval(value string) Option {
	return func(s *Config) {
		val, err := time.ParseDuration(value)
		if err != nil {
			log.Print(err.Error())
		}
		s.Store.StoreInterval = val
	}
}

func WithKey(value string) Option {
	return func(s *Config) {
		s.Server.Key = &value
	}
}

func WithDatabase(value string) Option {
	return func(s *Config) {
		s.Store.DatabaseConnectionString = &value
	}
}
