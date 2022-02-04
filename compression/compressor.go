package compression

import (
	"errors"
	"hzip/input"
	"hzip/output"
)

type Compressor struct {
	Output          output.Output
	Inputs          []input.Input
	frequency_table *FrequencyTable
}

func (compressor *Compressor) Compress() error {
	return errors.New("not implemented")
}

func (compressor *Compressor) Dump() error {
	return errors.New("not implemented")
}

func (compressor *Compressor) AddInput(input_obj input.Input) {
	compressor.Inputs = append(compressor.Inputs, input_obj)
}

func (compressor *Compressor) SetOutput(output_obj output.Output) {
	compressor.Output = output_obj
}
