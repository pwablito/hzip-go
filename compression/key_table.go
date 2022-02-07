package compression

import "errors"

type KeyTableData struct {
	length int // bits
	data   string
}

type KeyTable struct {
	table map[string]KeyTableData
}

func (table *KeyTable) Add(key string, data string, length int) {
	table.table[key] = KeyTableData{
		length: length,
		data:   data,
	}
}

func (table *KeyTable) ReadTree(tree *HuffmanTree) error {
	return errors.New("[ERROR] KeyTable.ReadTree(HuffmanTree) not implemented")
}
