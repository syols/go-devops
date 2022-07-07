package store

import (
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

func (f FileStore) Save(value []models.Payload) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := os.Remove(f.storeFile); err != nil {
		log.Print(err.Error())
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
			log.Print(err.Error())
		}
	}(file)

	return nil
}

func (f FileStore) Load() ([]models.Payload, error) {
	file, err := ioutil.ReadFile(f.storeFile)
	if err != nil {
		return nil, err
	}

	var payload []models.Payload
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
