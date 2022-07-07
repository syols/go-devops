package store

import (
	"encoding/json"
	"github.com/syols/go-devops/internal/metric"
	"io/ioutil"
	"log"
	"os"
)

type FileStore struct {
	storeFile string
}

func NewFileStore(storeFile string) FileStore {
	log.Printf("Storage from %s", storeFile)
	return FileStore{
		storeFile: storeFile,
	}
}

func (f FileStore) Save(value []metric.Payload) error {
	if len(value) == 0 {
		return nil
	}

	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := os.Remove(f.storeFile); err != nil {
		log.Printf("file %s not exist", f.storeFile)
	}

	file, err := os.Create(f.storeFile)
	if err != nil {
		return err
	}

	if _, err := file.Write(jsonBytes); err != nil {
		return err
	}

	defer func(file *os.File) {
		log.Printf("Save metrics to %s: %s", f.storeFile, string(jsonBytes))
		err := file.Close()
		if err != nil {
			log.Print(err.Error())
		}
	}(file)

	return nil
}

func (f FileStore) Load() ([]metric.Payload, error) {
	log.Printf("Load metrics from %s", f.storeFile)
	file, err := ioutil.ReadFile(f.storeFile)
	if err != nil {
		return nil, err
	}
	var payload []metric.Payload
	err = json.Unmarshal([]byte(file), &payload)
	return payload, err
}

func (f FileStore) Check() error {
	_, err := os.Stat(f.storeFile)
	return err
}

func (f FileStore) Type() string {
	return "file"
}
