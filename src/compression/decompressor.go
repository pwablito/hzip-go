package compression

import (
	"bytes"
	"errors"
	"fmt"
	"hzip/src/huffman_tree"
	"hzip/src/key_table"
	"os"

	"github.com/dgryski/go-bitstream"
)

type Decompressor struct {
	InputFilename string
	keytable      key_table.KeyTable
	reader        *bitstream.BitReader
	tree          *huffman_tree.HuffmanTree
}

func (decompressor *Decompressor) ReadMeta() error {
	file, err := os.Open(decompressor.InputFilename)
	if err != nil {
		return errors.New("Couldn't open archive: " + decompressor.InputFilename)
	}
	bits_read := 0
	decompressor.reader = bitstream.NewReader(file)
	num_table_entries, err := decompressor.reader.ReadBits(64)
	bits_read += 64
	if err != nil {
		fmt.Println("[ERROR] Couldn't read table size")
	}
	for i := 0; i < int(num_table_entries); i++ {
		key, err := decompressor.reader.ReadByte()
		bits_read += 8
		if err != nil {
			return errors.New("[ERROR] Couldn't read key")
		}
		val_length, err := decompressor.reader.ReadBits(64)
		bits_read += 64
		if err != nil {
			return errors.New("[ERROR] Couldn't read length")
		}
		var val_buffer bytes.Buffer
		val_buffer_writer := bitstream.NewWriter(&val_buffer)
		for j := 0; j < int(val_length); j++ {
			current_bit, err := decompressor.reader.ReadBit()
			bits_read += 1
			if err != nil {
				return errors.New("[ERROR] Couldn't read value")
			}
			val_buffer_writer.WriteBit(current_bit)
		}
		val_buffer_writer.Flush(bitstream.Zero)
		decompressor.keytable.Add(key, val_buffer, int(val_length))
	}
	// Flush out the padding bits
	decompressor.reader.ReadBits(8 - (bits_read % 8))
	// Now we have the key table, we can convert it to a huffman tree for fast decompression lookups
	decompressor.tree, err = decompressor.keytable.WriteTree()
	if err != nil {
		fmt.Println(err)
		return errors.New("[ERROR] Failed to generate huffman tree from key table")
	}
	return nil
}

func (decompressor Decompressor) CreateDirectoryStructure() error {
	// TODO possibly should collect directory structure in ReadMeta
	return errors.New("[ERROR] Not implemented")
}

func (decompressor Decompressor) Decompress() error {
	return errors.New("[ERROR] Not implemented")
}
