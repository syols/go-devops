package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/syols/go-devops/internal/metric"
	"log"
)

var schema = `CREATE TABLE metrics (
				  id VARCHAR(256) CONSTRAINT metric_id_pk PRIMARY KEY NOT NULL,
				  metric_type VARCHAR(256) NOT NULL,
				  counter_value NUMERIC DEFAULT NULL,
				  gauge_value DOUBLE PRECISION DEFAULT NULL,
				  hash VARCHAR(256) DEFAULT NULL);`

var insert = `INSERT INTO metrics (id, metric_type, counter_value, gauge_value, hash)
			  VALUES (:id, :metric_type, :counter_value, :gauge_value, :hash)
			  ON CONFLICT (id) DO UPDATE
			  SET metric_type = excluded.metric_type,
					counter_value = excluded.counter_value,
					gauge_value = excluded.gauge_value,
					hash = excluded.hash;`

var value = `SELECT id, metric_type, counter_value, gauge_value, hash FROM metrics`

type DatabaseStore struct {
	connectionString string
}

func NewDatabaseStore(connectionString string) DatabaseStore {
	log.Printf("Storage from %s", connectionString)
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if _, err := db.Exec(schema); err != nil {
		log.Printf(err.Error())
	}

	return DatabaseStore{
		connectionString: connectionString,
	}
}

func (d DatabaseStore) Save(value []metric.Payload) error {
	if len(value) == 0 {
		return nil
	}

	db, err := sqlx.Connect("postgres", d.connectionString)
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(db)
	if err != nil {
		return err
	}

	if _, err = db.NamedExec(insert, value); err != nil {
		return err
	}
	return nil
}

func (d DatabaseStore) Load() ([]metric.Payload, error) {
	var payload []metric.Payload
	db, err := sqlx.Connect("postgres", d.connectionString)
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(db)
	if err != nil {
		return payload, err
	}

	if err := db.Select(&payload, value); err != nil {
		return payload, err
	}
	return payload, nil
}

func (d DatabaseStore) Check() error {
	db, err := sqlx.Connect("postgres", d.connectionString)
	if err == nil {
		defer db.Close()
	}
	return err
}

func (d DatabaseStore) Type() string {
	return "database"
}
