package store

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/syols/go-devops/internal/metric"
	"log"
)

type DatabaseStore struct {
	connectionString string
}

func NewDatabaseStore(connectionString string) DatabaseStore {
	log.Printf("Storage from %s", connectionString)

	return DatabaseStore{
		connectionString: connectionString,
	}
}

func (d DatabaseStore) Save(value []metric.Payload) error {
	return nil
}

func (d DatabaseStore) Load() ([]metric.Payload, error) {
	var payload []metric.Payload
	return payload, nil
}

func (d DatabaseStore) Check() error {
	connection, err := pgx.Connect(context.Background(), d.connectionString)
	if err == nil {
		defer connection.Close(context.Background())
	}
	return err
}
