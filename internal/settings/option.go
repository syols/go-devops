package settings

import (
	"log"
	"net"
	"strconv"
	"time"
)

type Option func(s *Config)

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
