package settings

import (
	"log"
	"net"
	"strconv"
	"time"
)

type Option func(s *Settings)

func WithAddress(address string) Option {
	return func(s *Settings) {
		if host, port, err := net.SplitHostPort(address); err == nil {
			if port, err := strconv.ParseUint(port, 0, 16); err == nil {
				log.Printf("Address:\t%s", address)
				s.Server.Address.Host = host
				s.Server.Address.Port = uint16(port)
			}
		}
	}
}

func WithRestore(value string) Option {
	return func(s *Settings) {
		val, err := strconv.ParseBool(value)
		if err != nil {
			log.Printf("incorrect option WithRestore")
		}
		log.Printf("Is restore:\t%s", value)
		s.Store.Restore = val
	}
}

func WithStoreFile(value string) Option {
	return func(s *Settings) {
		log.Printf("StoreFile: %s", value)
		s.Store.StoreFile = &value
	}
}

func WithPollInterval(value string) Option {
	return func(s *Settings) {
		val, err := time.ParseDuration(value)
		if err != nil {
			log.Printf("incorrect option 'WithPollInterval'")
		}
		log.Printf("Poll interval:\t%s", value)
		s.Agent.PollInterval = val
	}
}

func WithReportInterval(value string) Option {
	return func(s *Settings) {
		val, err := time.ParseDuration(value)
		if err != nil {
			log.Printf("incorrect option 'WithReportInterval'")
		}
		log.Printf("Report interval:\t%s", value)
		s.Agent.ReportInterval = val
	}
}

func WithClientTimeout(value string) Option {
	return func(s *Settings) {
		val, err := time.ParseDuration(value)
		if err != nil {
			log.Printf("incorrect option 'WithClientTimeout'")
		}
		log.Printf("Client timeout interval:\t%s", value)
		s.Agent.ClientTimeout = val
	}
}

func WithStoreInterval(value string) Option {
	return func(s *Settings) {
		val, err := time.ParseDuration(value)
		if err != nil {
			log.Printf("incorrect option 'WithStoreInterval'")
		}
		log.Printf("Store interval: %s", value)
		s.Store.StoreInterval = val
	}
}

func WithKey(value string) Option {
	return func(s *Settings) {
		log.Printf("Sha256 key:\t%s", value)
		s.Server.Key = &value
	}
}

func WithDatabase(value string) Option {
	return func(s *Settings) {
		log.Printf("Database:\t%s", value)
		s.Store.DatabaseConnectionString = &value
	}
}
