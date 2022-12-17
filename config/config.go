package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Option function of a certain type
type Option func(s *Config)

// Config struct
type Config struct {
	Server ServerConfig `yaml:"server" json:"server"`
	Agent  AgentConfig  `yaml:"agent" json:"agent"`
	Store  StoreConfig  `yaml:"store" json:"store"`
	Grpc   *GrpcConfig  `yaml:"grpc,omitempty" json:"grpc,omitempty"`
}

// ServerConfig Server config struct
type ServerConfig struct {
	Address       Address `yaml:"address" json:"address"`
	Key           *string `yaml:"key,omitempty" json:"key,omitempty"`
	TrustedSubnet *string `yaml:"trusted_subnet,omitempty" json:"trusted_subnet,omitempty"`
}

// AgentConfig Agent config struct
type AgentConfig struct {
	PollInterval   time.Duration `yaml:"poll_interval" json:"poll_interval"`
	ReportInterval time.Duration `yaml:"report_interval" json:"report_interval"`
	ClientTimeout  time.Duration `yaml:"client_timeout" json:"client_timeout"`
}

// StoreConfig Store config struct
type StoreConfig struct {
	CryptoKeyFilePath        *string
	DatabaseConnectionString *string       `yaml:"database,omitempty" json:"database,omitempty"`
	StoreFile                *string       `yaml:"store_file,omitempty" json:"store_file,omitempty"`
	Restore                  bool          `yaml:"restore" json:"restore"`
	StoreInterval            time.Duration `yaml:"store_interval" json:"store_interval"`
}

// GrpcConfig Store Grpc config
type GrpcConfig struct {
	Address Address `yaml:"address" json:"address"`
}

// Address struct
type Address struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

// NewConfig creates config struct
func NewConfig() (settings Config) {
	err := settings.setDefault("develop.json")
	if err != nil {
		return Config{}
	}
	settings.LoadFromEnvironment()
	return settings
}

// LoadFromEnvironment config struct
func (s *Config) LoadFromEnvironment() {
	s.SetFromOptions(NewEnvironmentVariables().Options()...)
}

// Address create HTTP address
func (a *Address) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

// String create string from config
func (s *Config) String() (result string) {
	if marshal, err := yaml.Marshal(s); err != nil {
		result = string(marshal)
	}
	return
}

func (s *Config) SetFromOptions(options ...Option) {
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

func withHTTPAddress(address string) Option {
	return func(s *Config) {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			log.Fatal(err.Error())
		}

		value, err := strconv.ParseUint(port, 0, 16)
		if err != nil {
			log.Fatal(err.Error())
		}

		s.Server.Address.Host = host
		s.Server.Address.Port = uint16(value)
	}
}

func withGrpcAddress(address string) Option {
	return func(s *Config) {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			log.Fatal(err.Error())
		}

		value, err := strconv.ParseUint(port, 0, 16)
		if err != nil {
			log.Fatal(err.Error())
		}

		s.Grpc = &GrpcConfig{
			Address: Address{
				Host: host,
				Port: uint16(value),
			},
		}
	}
}

func withTrustedSubnet(value string) Option {
	return func(s *Config) {
		s.Server.TrustedSubnet = &value
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

func ReplaceNoneValue(value string) string {
	if len(value) == 0 {
		return "N/A"
	}
	return value
}
