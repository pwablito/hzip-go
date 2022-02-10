package frequency_table

type FrequencyTable struct {
	frequencies map[byte]int
}

func (freq_table *FrequencyTable) Increment(key byte) {
	_, present := freq_table.frequencies[key]
	if present {
		freq_table.frequencies[key] += 1
	} else {
		freq_table.frequencies[key] = 1
	}
}

func (freq_table *FrequencyTable) GetFrequencies() map[byte]int {
	return freq_table.frequencies
}
