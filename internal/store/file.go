package store

import (
	"context"
	"encoding/json"
	"github.com/syols/go-devops/internal/models"
	"io/ioutil"
	"log"
	"os"
)

type FileStore struct {
	storeFile string
}

func NewFileStore(storeFile string) FileStore {
	return FileStore{
		storeFile: storeFile,
	}
}

func (f FileStore) Save(ctx context.Context, value []models.Metric) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := os.Remove(f.storeFile); err != nil {
		return err
	}

	file, err := os.Create(f.storeFile)
	if err != nil {
		return err
	}

	if _, err := file.Write(jsonBytes); err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	return nil
}

func (f FileStore) Load(_ context.Context) ([]models.Metric, error) {
	file, err := ioutil.ReadFile(f.storeFile)
	if err != nil {
		return nil, err
	}

	var payload []models.Metric
	err = json.Unmarshal(file, &payload)

	return payload, err
}

func (f FileStore) Check() error {
	_, err := os.Stat(f.storeFile)
	return err
}

func (f FileStore) Type() string {
	return "file"
}
