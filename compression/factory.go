package compression

import "hzip/input"

func CreateCompressor() Compressor {
	return Compressor{
		Inputs:          make([]input.Input, 0),
		Output:          nil,
		frequency_table: nil,
	}
}
