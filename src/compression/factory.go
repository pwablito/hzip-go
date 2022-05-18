package compression

import (
	"hzip/src/input"
	"hzip/src/key_table"
)

func CreateCompressor() Compressor {
	return Compressor{
		Inputs: make([]input.Input, 0),
		Output: nil,
	}
}

func CreateDecompressor(filename string) Decompressor {
	return Decompressor{
		InputFilename: filename,
		keyTable:      key_table.CreateKeyTable(),
		reader:        nil,
	}
}
