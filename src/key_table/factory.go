package key_table

func CreateKeyTable() KeyTable {
	return KeyTable{
		table: make(map[byte]KeyTableData),
	}
}
