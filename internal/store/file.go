package store

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/syols/go-devops/internal/models"
)

// FileStore struct
type FileStore struct {
	storeFile string
}

// NewFileStore creates file store
func NewFileStore(storeFile string) FileStore {
	return FileStore{
		storeFile: storeFile,
	}
}

// Save to file store
func (f FileStore) Save(_ context.Context, value []models.Metric) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := os.Remove(f.storeFile); err != nil {
		log.Println(err)
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

// Load to file store
func (f FileStore) Load(_ context.Context) ([]models.Metric, error) {
	file, err := ioutil.ReadFile(f.storeFile)
	if err != nil {
		return nil, err
	}

	var payload []models.Metric
	err = json.Unmarshal(file, &payload)

	return payload, err
}

// Check in store
func (f FileStore) Check() error {
	_, err := os.Stat(f.storeFile)
	return err
}

// Type of store
func (f FileStore) Type() string {
	return "file"
}
