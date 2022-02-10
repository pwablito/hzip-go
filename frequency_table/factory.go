package frequency_table

func CreateFrequencyTable() FrequencyTable {
	return FrequencyTable{
		frequencies: make(map[byte]int),
	}
}
