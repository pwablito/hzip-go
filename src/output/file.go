package output

import (
	"errors"
	"fmt"
	"os"
)

type FileOutput struct {
	Filename string
	Mode     int
}

func (file_output *FileOutput) Write(data []byte) error {
	file, err := os.OpenFile(file_output.Filename, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return errors.New("[ERROR] Couldn't open file")
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return errors.New("[ERROR] Failed to write to file " + file_output.Filename)
	}
	err = file.Sync()
	if err != nil {
		return errors.New("[ERROR] Failed to synce file " + file_output.Filename)
	}
	return nil
}

func (file_output *FileOutput) Open() error {
	err := file_output.DeleteFileIfExists()
	if err != nil {
		fmt.Println(err)
		return errors.New("[ERROR] Failed to remove existing file")
	}
	err = file_output.CreateFile()
	if err != nil {
		fmt.Println(err)
		return errors.New("[ERROR] Failed to create file " + file_output.Filename)
	}
	return nil
}

func (file_output *FileOutput) Close() error {
	return nil
}

func (file_output *FileOutput) DeleteFileIfExists() error {
	_, err := os.Stat(file_output.Filename)
	if os.IsNotExist(err) {
		return nil
	}
	fmt.Println("[INFO] Removing " + file_output.Filename)
	err = os.Remove(file_output.Filename)
	if err != nil {
		return errors.New("[ERROR] Couldn't delete file " + file_output.Filename)
	}
	return nil
}

func (file_output *FileOutput) CreateFile() error {
	fmt.Println("[INFO] Creating " + file_output.Filename)
	file, err := os.Create(file_output.Filename)
	if err != nil {
		return errors.New("[ERROR] Couldn't create file")
	}
	file.Close()
	return nil
}
