package compression

import "errors"

type Decompressor struct {
	InputFilename string
}

func (decompressor Decompressor) CreateDirectoryStructure() error {
	return errors.New("[ERROR] Not implemented")
}

func (decompressor Decompressor) Decompress() error {
	return errors.New("[ERROR] Not implemented")
}
