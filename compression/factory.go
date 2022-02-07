package compression

import "hzip/input"

func CreateCompressor() Compressor {
	return Compressor{
		Inputs:          make([]input.Input, 0),
		Output:          nil,
		frequency_table: CreateFrequencyTable(),
	}
}
func CreateFrequencyTable() FrequencyTable {
	return FrequencyTable{
		frequencies: make(map[byte]int),
	}
}

func CreateKeyTable() KeyTable {
	return KeyTable{
		table: make(map[string]KeyTableData),
	}
}
