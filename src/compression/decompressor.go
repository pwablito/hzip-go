package compression

import (
	"bytes"
	"errors"
	"fmt"
	"hzip/src/huffman_tree"
	"hzip/src/key_table"
	"os"
	"path/filepath"

	"github.com/dgryski/go-bitstream"
	"github.com/schollz/progressbar/v3"
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
			bits_read++
			if err != nil {
				return errors.New("[ERROR] Couldn't read value")
			}
			val_buffer_writer.WriteBit(current_bit)
		}
		val_buffer_writer.Flush(bitstream.Zero)
		decompressor.keytable.Add(key, val_buffer, int(val_length))
	}
	// Flush out the padding bits
	if bits_read%8 != 0 {
		decompressor.reader.ReadBits(8 - (bits_read % 8))
	}
	// Now we have the key table, we can convert it to a huffman tree for fast decompression lookups
	decompressor.tree, err = decompressor.keytable.WriteTree()
	if err != nil {
		fmt.Println(err)
		return errors.New("[ERROR] Failed to generate huffman tree from key table")
	}
	return nil
}

func (decompressor Decompressor) Decompress() error {
	// TODO possibly should collect directory structure in ReadMeta
	reader := decompressor.reader
	num_files, err := reader.ReadBits(64)
	if err != nil {
		return errors.New("[ERROR] Couldn't get number of files")
	}
	bar := progressbar.NewOptions(
		int(num_files),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetPredictTime(true),
	)
	for i := 0; i < int(num_files); i++ {
		bar.Add(1)
		bits_read := 0
		filename_len, err := reader.ReadBits(64)
		bits_read += 64
		if err != nil {
			return errors.New("[ERROR] Couldn't read filename length")
		}
		var filename_buffer bytes.Buffer
		filename_writer := bitstream.NewWriter(&filename_buffer)
		for j := 0; j < int(filename_len); j++ {
			byte_obj, err := reader.ReadByte()
			bits_read += 8
			if err != nil {
				return errors.New("[ERROR] Couldn't read filename")
			}
			err = filename_writer.WriteByte(byte_obj)
			if err != nil {
				return errors.New("[ERROR] Couldn't write filename to buffer")
			}
		}

		// Read past the file contents so we can cleanly get the next filename
		buffer_len, err := reader.ReadBits(64)
		bits_read += 64
		if err != nil {
			return errors.New("[ERROR] Couldn't read file length")
		}
		// initialize current_chunk
		var current_chunk bytes.Buffer = *bytes.NewBuffer(make([]byte, 0))
		current_chunk_writer := bitstream.NewWriter(&current_chunk)
		current_chunk_len := 0
		var decompressed_buffer bytes.Buffer = *bytes.NewBuffer(make([]byte, 0))
		decompressed_buffer_writer := bitstream.NewWriter(&decompressed_buffer)
		for j := 0; j < int(buffer_len); j++ {
			// Get a bit from the compressed buffer
			bit, err := reader.ReadBit()
			bits_read++
			if err != nil {
				return errors.New("[ERROR] Couldn't seek past file buffer")
			}
			// Write bit to current chunk
			current_chunk_writer.WriteBit(bit)
			// Increment length
			current_chunk_len++
			// Fill the rest with 0s
			current_chunk_writer.Flush(bitstream.Zero)
			// Perform htree lookup
			data, found_leaf, err := decompressor.tree.Lookup(current_chunk, current_chunk_len)
			if err != nil {
				fmt.Println(err)
				return errors.New("[ERROR] Overflow occurred")
			}
			if !found_leaf {
				// Create a new buffer
				var new_buffer bytes.Buffer = *bytes.NewBuffer(make([]byte, 0))
				new_writer := bitstream.NewWriter(&new_buffer)
				current_chunk_reader := bitstream.NewReader(&current_chunk)
				// Copy current_chunk to new_buffer with no padding
				for k := 0; k < current_chunk_len; k++ {
					bit, err = current_chunk_reader.ReadBit()
					if err != nil {
						return errors.New("[ERROR] Couldn't transfer data")
					}
					new_writer.WriteBit(bit)
				}
				current_chunk = new_buffer
				current_chunk_writer = bitstream.NewWriter(&current_chunk)
			} else {
				decompressed_buffer_writer.WriteByte(data)
				current_chunk = *bytes.NewBuffer(make([]byte, 0))
				current_chunk_writer = bitstream.NewWriter(&current_chunk)
			}
		}
		// Create and write file
		dirpath := filepath.Dir(filename_buffer.String()) // split here
		err = os.MkdirAll(dirpath, 0o755)                 // TODO track modes in archive
		if err != nil {
			return errors.New("[ERROR] Couldn't create directory " + dirpath)
		}
		file, err := os.Create(filename_buffer.String())
		if err != nil {
			return errors.New("[ERROR] Couldn't open file " + filename_buffer.String())
		}
		file.Write(decompressed_buffer.Bytes())
		file.Close()

		// Reset byte boundary
		if bits_read%8 != 0 {
			bits, err := reader.ReadBits(8 - (bits_read % 8))
			if bits > 0 {
				return errors.New("[ERROR] Expected bits to be zero")
			}
			if err != nil {
				return errors.New("[ERROR] Couldn't zero out buffer")
			}
		}
	}
	bar.Finish()
	return nil
}
