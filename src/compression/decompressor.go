package compression

import (
	"errors"
	"hzip/src/key_table"
)

type Decompressor struct {
	InputFilename string
	keytable      *key_table.KeyTable
}

func (decompressor *Decompressor) ReadMeta() error {
	return errors.New("[ERROR] Not implemented")
}

func (decompressor Decompressor) CreateDirectoryStructure() error {
	return errors.New("[ERROR] Not implemented")
}

func (decompressor Decompressor) Decompress() error {
	return errors.New("[ERROR] Not implemented")
}
