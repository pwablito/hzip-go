package compression

import (
	"errors"
	"hzip/input"
	"hzip/output"
)

type Compressor struct {
	Output output.Output
	Inputs []input.Input
}

func (compressor *Compressor) Compress() error {
	return errors.New("not implemented")
}

func (compressor *Compressor) Dump() error {
	return errors.New("not implemented")
}
