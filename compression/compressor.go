package compression

import (
	"errors"
	"fmt"
	"hzip/input"
	"hzip/output"
)

type Compressor struct {
	Output          output.Output
	Inputs          []input.Input
	frequency_table FrequencyTable
}

func (compressor *Compressor) Compress() error {
	fmt.Println("[INFO] Creating frequency table")
	for _, input_obj := range compressor.Inputs {
		data, err := input_obj.GetData()
		if err != nil {
			return errors.New("[ERROR] Failed to read data from input")
		}
		for _, current_byte := range data {
			compressor.frequency_table.Increment(current_byte)
		}
	}
	return errors.New("[ERROR] not fully implemented")
}

func (compressor *Compressor) Dump() error {
	return errors.New("[ERROR] Not implemented")
}

func (compressor *Compressor) AddInput(input_obj input.Input) {
	compressor.Inputs = append(compressor.Inputs, input_obj)
}

func (compressor *Compressor) SetOutput(output_obj output.Output) {
	compressor.Output = output_obj
}
