package compression

import (
	"hzip/frequency_table"
	"hzip/input"
)

func CreateCompressor() Compressor {
	return Compressor{
		Inputs:          make([]input.Input, 0),
		Output:          nil,
		frequency_table: frequency_table.CreateFrequencyTable(),
	}
}
