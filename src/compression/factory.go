package compression

import (
	"hzip/src/frequency_table"
	"hzip/src/input"
)

func CreateCompressor() Compressor {
	return Compressor{
		Inputs:          make([]input.Input, 0),
		Output:          nil,
		frequency_table: frequency_table.CreateFrequencyTable(),
	}
}
