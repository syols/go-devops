package internal

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type Settings struct {
	Server ServerSettings `yaml:"server"`
	Agent  Agent          `yaml:"agent"`
}

type ServerSettings struct {
	Address       Address       `yaml:"address"`
	Restore       bool          `yaml:"restore"`
	StoreInterval time.Duration `yaml:"store_interval"`
	StoreFile     string        `yaml:"store_file"`
}

type Address struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

type Agent struct {
	PollInterval   time.Duration `yaml:"poll_interval"`
	ReportInterval time.Duration `yaml:"report_interval"`
	ClientTimeout  time.Duration `yaml:"client_timeout"`
}

type Option func(s *Settings)

func NewSettings(options ...Option) (settings Settings) {
	settings.setFromFile("configs/default.yml")
	settings.setFromEnv()
	settings.setFromOptions(options...)
	return settings
}

func WithHost(host string) Option {
	return func(s *Settings) {
		s.Server.Address.Host = host
	}
}

func WithPort(port uint16) Option {
	return func(s *Settings) {
		s.Server.Address.Port = port
	}
}

func WithRestore(restore bool) Option {
	return func(s *Settings) {
		s.Server.Restore = restore
	}
}

func WithStoreFile(fileName string) Option {
	return func(s *Settings) {
		s.Server.StoreFile = fileName
	}
}

func WithPollInterval(seconds int) Option {
	return func(s *Settings) {
		s.Agent.PollInterval = time.Duration(seconds) * time.Second
	}
}

func WithReportInterval(seconds int) Option {
	return func(s *Settings) {
		s.Agent.ReportInterval = time.Duration(seconds) * time.Second
	}
}

func WithClientTimeout(seconds int) Option {
	return func(s *Settings) {
		s.Agent.ClientTimeout = time.Duration(seconds) * time.Second
	}
}

func WithStoreInterval(seconds int) Option {
	return func(s *Settings) {
		s.Server.StoreInterval = time.Duration(seconds) * time.Second
	}
}

func (s *Settings) setFromOptions(options ...Option) {
	for _, fn := range options {
		fn(s)
	}
}

func (s *Settings) setFromFile(configPath string) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Config load error")
	}

	if err := yaml.Unmarshal(file, s); err != nil {
		log.Fatalf("Config load error")
	}

	s.Agent.PollInterval *= time.Second
	s.Agent.ReportInterval *= time.Second
	s.Agent.ClientTimeout *= time.Second
}

func (s *Settings) setFromEnv() {
	address := os.Getenv("ADDRESS")
	if host, port, err := net.SplitHostPort(address); err == nil {
		if port, err := strconv.ParseUint(port, 0, 16); err == nil {
			s.setFromOptions(
				WithHost(host),
				WithPort(uint16(port)),
			)
		}
	}

	if val, err := strconv.Atoi(os.Getenv("REPORT_INTERVAL")); err == nil {
		s.setFromOptions(WithReportInterval(val))
	}

	if val, err := strconv.Atoi(os.Getenv("POLL_INTERVAL")); err == nil {
		s.setFromOptions(WithPollInterval(val))
	}

	if val, err := strconv.Atoi(os.Getenv("CLIENT_TIMEOUT")); err == nil {
		s.setFromOptions(WithClientTimeout(val))
	}

	if val, err := strconv.Atoi(os.Getenv("STORE_INTERVAL")); err == nil {
		s.setFromOptions(WithStoreInterval(val))
	}

	if val, err := strconv.ParseBool(os.Getenv("RESTORE")); err == nil {
		s.setFromOptions(WithRestore(val))
	}

	s.setFromOptions(WithStoreFile(os.Getenv("STORE_FILE")))
}

func (s *Settings) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Server.Address.Host, s.Server.Address.Port)
}
