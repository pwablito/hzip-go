package output

import "errors"

type FileOutput struct {
	Filename string
}

func (file_output FileOutput) Write(data []byte) error {
	return errors.New("not implemented")
}
