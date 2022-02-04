package compression

type FrequencyTable struct {
	frequencies map[string]int
}

func (freq_table *FrequencyTable) Increment(key []byte) {
	str_key := string(key)
	_, present := freq_table.frequencies[str_key]
	if present {
		freq_table.frequencies[str_key] += 1
	} else {
		freq_table.frequencies[str_key] = 1
	}
}
