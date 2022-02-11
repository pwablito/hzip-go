package output

import (
	"errors"
	"fmt"
	"os"
)

type FileOutput struct {
	Filename string
	Mode     int
	file     *os.File
}

func (file_output *FileOutput) Write(data []byte) error {
	_, err := file_output.file.Write(data)
	if err != nil {
		return errors.New("[ERROR] Failed to write to file " + file_output.Filename)
	}
	return nil
}

func (file_output *FileOutput) Open() error {
	err := file_output.CreateFileIfNotExists()
	if err != nil {
		fmt.Println(err)
		return errors.New("[ERROR] Failed to create file " + file_output.Filename)
	}
	file_output.file, err = os.Open(file_output.Filename)
	if err != nil {
		return errors.New("[ERROR] Failed to open file after creation: " + file_output.Filename)
	}
	return nil
}

func (file_output *FileOutput) Close() error {
	err := file_output.file.Sync()
	if err != nil {
		return errors.New("[ERROR] Failed to sync file " + file_output.Filename)
	}
	file_output.file.Close()
	return nil
}

func (file_output *FileOutput) CreateFileIfNotExists() error {
	_, err := os.Stat(file_output.Filename)
	if os.IsNotExist(err) {
		fmt.Println("[INFO] Creating " + file_output.Filename)
		file_output.file, err = os.Create(file_output.Filename)
		if err != nil {
			return errors.New("[ERROR] Couldn't create file")
		}
		defer file_output.file.Close()
	}
	return nil
}
