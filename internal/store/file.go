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

	if removeErr := os.Remove(f.storeFile); removeErr != nil {
		log.Println(removeErr)
	}

	file, creatErr := os.Create(f.storeFile)
	if creatErr != nil {
		return creatErr
	}

	if _, writErr := file.Write(jsonBytes); writErr != nil {
		return writErr
	}

	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			log.Fatal(closeErr)
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
