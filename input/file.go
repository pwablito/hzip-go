package input

import (
	"errors"
	"os"
)

type FileInput struct {
	Filename string
}

func (file_input FileInput) Read() ([]byte, error) {
	data, err := os.ReadFile(file_input.Filename)
	if err != nil {
		return nil, errors.New("failed to read file")
	}
	return data, nil
}
