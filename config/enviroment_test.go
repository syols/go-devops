package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	duration := time.Duration(1)
	address := "0.0.0.0:8080"
	reportInterval := duration.String()
	pollInterval := duration.String()
	clientTimeout := duration.String()
	storeInterval := duration.String()
	restore := "true"
	key := "some_key"
	storeFile := "/tmp/devops-metrics-db.json"
	databaseDsn := "postgres://postgres:postgres@localhost/postgres?sslmode=disable"

	t.Setenv("ADDRESS", address)
	t.Setenv("REPORT_INTERVAL", reportInterval)
	t.Setenv("POLL_INTERVAL", pollInterval)
	t.Setenv("CLIENT_TIMEOUT", clientTimeout)
	t.Setenv("STORE_INTERVAL", storeInterval)
	t.Setenv("RESTORE", restore)
	t.Setenv("KEY", key)
	t.Setenv("STORE_FILE", storeFile)
	t.Setenv("DATABASE_DSN", databaseDsn)

	config := Config{}
	options := NewEnvironmentVariables().Options()
	config.SetFromOptions(options...)
	assert.Equal(t, address, config.Server.Address.String())
	assert.Equal(t, duration, config.Agent.ReportInterval)
	assert.Equal(t, duration, config.Agent.PollInterval)
	assert.Equal(t, duration, config.Agent.ClientTimeout)
	assert.Equal(t, duration, config.Store.StoreInterval)
	assert.Equal(t, true, config.Store.Restore)
	assert.Equal(t, key, *config.Server.Key)
	assert.Equal(t, storeFile, *config.Store.StoreFile)
	assert.Equal(t, databaseDsn, *config.Store.DatabaseConnectionString)
}
