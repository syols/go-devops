package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type Option func(s *Config)

type Config struct {
	Server ServerConfig `yaml:"server" json:"server"`
	Agent  AgentConfig  `yaml:"agent" json:"agent"`
	Store  StoreConfig  `yaml:"store" json:"store"`
}

type ServerConfig struct {
	Address Address `yaml:"address" json:"address"`
	Key     *string `yaml:"key,omitempty" json:"key,omitempty"`
}

type AgentConfig struct {
	PollInterval   time.Duration `yaml:"poll_interval" json:"poll_interval"`
	ReportInterval time.Duration `yaml:"report_interval" json:"report_interval"`
	ClientTimeout  time.Duration `yaml:"client_timeout" json:"client_timeout"`
}

type StoreConfig struct {
	CryptoKeyFilePath        *string
	DatabaseConnectionString *string       `yaml:"database,omitempty" json:"database,omitempty"`
	StoreFile                *string       `yaml:"store_file,omitempty" json:"store_file,omitempty"`
	Restore                  bool          `yaml:"restore" json:"restore"`
	StoreInterval            time.Duration `yaml:"store_interval" json:"store_interval"`
}

type Address struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

func NewConfig() (settings Config) {
	err := settings.setDefault("develop.json")
	if err != nil {
		return Config{}
	}
	settings.setFromOptions(newVariables().Options()...)
	return settings
}

func (s *Config) setFromOptions(options ...Option) {
	for _, fn := range options {
		fn(s)
	}
}

func (s *Config) setDefault(configPath string) error {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, s); err != nil {
		return err
	}
	return nil
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

func withAddress(address string) Option {
	return func(s *Config) {
		if host, port, err := net.SplitHostPort(address); err == nil {
			if port, err := strconv.ParseUint(port, 0, 16); err == nil {
				s.Server.Address.Host = host
				s.Server.Address.Port = uint16(port)
			}
		}
	}
}

func withStoreFile(value string) Option {
	return func(s *Config) {
		s.Store.StoreFile = &value
	}
}

func withRestore(value string) Option {
	return func(s *Config) {
		if val, err := strconv.ParseBool(value); err == nil {
			s.Store.Restore = val
		}
	}
}

func withPollInterval(value string) Option {
	return func(s *Config) {
		if val, err := time.ParseDuration(value); err == nil {
			s.Agent.PollInterval = val
		}
	}
}

func withReportInterval(value string) Option {
	return func(s *Config) {
		if val, err := time.ParseDuration(value); err == nil {
			s.Agent.ReportInterval = val
		}
	}
}

func withClientTimeout(value string) Option {
	return func(s *Config) {
		if val, err := time.ParseDuration(value); err == nil {
			s.Agent.ClientTimeout = val
		}

	}
}

func withStoreInterval(value string) Option {
	return func(s *Config) {
		if val, err := time.ParseDuration(value); err == nil {
			s.Store.StoreInterval = val
		}

	}
}

func withKey(value string) Option {
	return func(s *Config) {
		s.Server.Key = &value
	}
}

func withDatabase(value string) Option {
	return func(s *Config) {
		s.Store.DatabaseConnectionString = &value
	}
}

func withCryptoKey(value string) Option {
	return func(s *Config) {
		s.Store.CryptoKeyFilePath = &value
	}
}
