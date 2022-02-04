package compression

import (
	"errors"
	"fmt"
	"hzip/input"
	"hzip/output"

	"github.com/schollz/progressbar/v3"
)

type Compressor struct {
	Output          output.Output
	Inputs          []input.Input
	frequency_table FrequencyTable
}

func (compressor *Compressor) Compress() error {
	fmt.Println("[INFO] Creating frequency table")
	bar := progressbar.NewOptions(
		len(compressor.Inputs),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetPredictTime(true),
	)
	for _, input_obj := range compressor.Inputs {
		bar.Add(1)
		data, err := input_obj.GetData()
		if err != nil {
			fmt.Println(err)
			return errors.New("[ERROR] Failed to read data from input")
		}
		for _, current_byte := range data {
			compressor.frequency_table.Increment(current_byte)
		}
	}
	bar.Finish()
	return errors.New("[ERROR] Compression not fully implemented")
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
