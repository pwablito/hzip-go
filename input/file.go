package input

import (
	"errors"
	"fmt"
	"os"
)

type FileInput struct {
	Filename string
}

func (file_input FileInput) GetData() ([]byte, error) {
	data, err := os.ReadFile(file_input.Filename)
	if err != nil {
		fmt.Println("[ERROR]", err)
		return nil, errors.New("failed to read file")
	}
	return data, nil
}
