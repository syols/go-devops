package store

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/syols/go-devops/internal/models"
)

// DatabaseStore struct
type DatabaseStore struct {
	dataSourceName string
	selectQuery    string
	insertQuery    string
}

// NewDatabaseStore creates database store
func NewDatabaseStore(connectionString string) (DatabaseStore, error) {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return DatabaseStore{}, err
	}

	if _, err := db.Exec(loadSQL("create.sql")); err != nil {
		log.Print(err)
	}

	return DatabaseStore{
		dataSourceName: connectionString,
		selectQuery:    loadSQL("select.sql"),
		insertQuery:    loadSQL("insert.sql"),
	}, nil
}

// Save metrics to database
func (d DatabaseStore) Save(ctx context.Context, value []models.Metric) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", d.dataSourceName)
	if err != nil {
		return err
	}

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	if _, err = db.NamedExec(d.insertQuery, value); err != nil {
		return err
	}

	return nil
}

// Load metrics to database
func (d DatabaseStore) Load(ctx context.Context) ([]models.Metric, error) {
	var payload []models.Metric
	db, err := sqlx.ConnectContext(ctx, "postgres", d.dataSourceName)
	if err != nil {
		return payload, err
	}

	defer func(db *sqlx.DB) {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}(db)

	if err := db.Select(&payload, d.selectQuery); err != nil {
		return payload, err
	}

	return payload, nil
}

// Check metrics in database
func (d DatabaseStore) Check() error {
	db, err := sqlx.Connect("postgres", d.dataSourceName)
	if err == nil {
		defer func(db *sqlx.DB) {
			err := db.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(db)
	}

	return err
}

// Type of store
func (d DatabaseStore) Type() string {
	return "database"
}

func loadSQL(file string) string {
	path := filepath.Join("internal", "scripts", file)

	c, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return string(c)
}
