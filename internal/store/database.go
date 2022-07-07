package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/syols/go-devops/internal/model"
	"io/ioutil"
	"log"
	"path/filepath"
)

type DatabaseStore struct {
	dataSourceName string
	selectQuery    string
	insertQuery    string
}

func loadSQL(file string) string {
	path := filepath.Join("scripts", file)

	c, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return string(c)
}

func NewDatabaseStore(connectionString string) DatabaseStore {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if _, err := db.Exec(loadSQL("create.sql")); err != nil {
		log.Print(err.Error())
	}

	return DatabaseStore{
		dataSourceName: connectionString,
		selectQuery:    loadSQL("select.sql"),
		insertQuery:    loadSQL("insert.sql"),
	}
}

func (d DatabaseStore) Save(value []model.Payload) error {
	db, err := sqlx.Connect("postgres", d.dataSourceName)
	if err != nil {
		return err
	}

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(db)

	if _, err = db.NamedExec(d.insertQuery, value); err != nil {
		return err
	}

	return nil
}

func (d DatabaseStore) Load() ([]model.Payload, error) {
	var payload []model.Payload
	db, err := sqlx.Connect("postgres", d.dataSourceName)
	if err != nil {
		return payload, err
	}

	defer func(db *sqlx.DB) {
		if err := db.Close(); err != nil {
			log.Fatalf(err.Error())
		}
	}(db)

	if err := db.Select(&payload, d.selectQuery); err != nil {
		return payload, err
	}

	return payload, nil
}

func (d DatabaseStore) Check() error {
	db, err := sqlx.Connect("postgres", d.dataSourceName)
	if err == nil {
		defer db.Close()
	}

	return err
}

func (d DatabaseStore) Type() string {
	return "database"
}
