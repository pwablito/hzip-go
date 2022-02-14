package key_table

func CreateKeyTable() KeyTable {
	return KeyTable{
		Table: make(map[byte]KeyTableData),
	}
}
