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
	keyTable      key_table.KeyTable
	reader        *bitstream.BitReader
	tree          *huffman_tree.HuffmanTree
}

func (decompressor *Decompressor) ReadMeta() error {
	file, err := os.Open(decompressor.InputFilename)
	if err != nil {
		return errors.New("Couldn't open archive: " + decompressor.InputFilename)
	}
	bitsRead := 0
	decompressor.reader = bitstream.NewReader(file)
	numTableEntries, err := decompressor.reader.ReadBits(64)
	bitsRead += 64
	if err != nil {
		fmt.Println("[ERROR] Couldn't read table size")
	}
	for i := 0; i < int(numTableEntries); i++ {
		key, err := decompressor.reader.ReadByte()
		bitsRead += 8
		if err != nil {
			return errors.New("[ERROR] Couldn't read key")
		}
		valLength, err := decompressor.reader.ReadBits(64)
		bitsRead += 64
		if err != nil {
			return errors.New("[ERROR] Couldn't read length")
		}
		var valBuffer bytes.Buffer
		valBufferWriter := bitstream.NewWriter(&valBuffer)
		for j := 0; j < int(valLength); j++ {
			currentBit, err := decompressor.reader.ReadBit()
			bitsRead++
			if err != nil {
				return errors.New("[ERROR] Couldn't read value")
			}
			err = valBufferWriter.WriteBit(currentBit)
			if err != nil {
				return errors.New("[ERROR] Failed to write bit to stream")
			}
		}
		err = valBufferWriter.Flush(bitstream.Zero)
		if err != nil {
			return errors.New("[ERROR] Failed to flush bitstream")
		}
		decompressor.keyTable.Add(key, valBuffer, int(valLength))
	}
	// Flush out the padding bits
	if bitsRead%8 != 0 {
		_, err := decompressor.reader.ReadBits(8 - (bitsRead % 8))
		if err != nil {
			return errors.New("[ERROR] Failed to flush bits by reading")
		}
	}
	// Now we have the key table, we can convert it to a huffman tree for fast decompression lookups
	decompressor.tree, err = decompressor.keyTable.WriteTree()
	if err != nil {
		fmt.Println(err)
		return errors.New("[ERROR] Failed to generate huffman tree from key table")
	}
	return nil
}

func (decompressor Decompressor) Decompress() error {
	// TODO possibly should collect directory structure in ReadMeta
	reader := decompressor.reader
	numFiles, err := reader.ReadBits(64)
	if err != nil {
		return errors.New("[ERROR] Couldn't get number of files")
	}
	bar := progressbar.NewOptions(
		int(numFiles),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetPredictTime(true),
	)
	for i := 0; i < int(numFiles); i++ {
		err := bar.Add(1)
		if err != nil {
			return errors.New("[ERROR] Failed to modify progress bar status")
		}
		bitsRead := 0
		filenameLen, err := reader.ReadBits(64)
		bitsRead += 64
		if err != nil {
			return errors.New("[ERROR] Couldn't read filename length")
		}
		var filenameBuffer bytes.Buffer
		filenameWriter := bitstream.NewWriter(&filenameBuffer)
		for j := 0; j < int(filenameLen); j++ {
			byteObj, err := reader.ReadByte()
			bitsRead += 8
			if err != nil {
				return errors.New("[ERROR] Couldn't read filename")
			}
			err = filenameWriter.WriteByte(byteObj)
			if err != nil {
				return errors.New("[ERROR] Couldn't write filename to buffer")
			}
		}

		// Read past the file contents, so we can cleanly get the next filename
		bufferLen, err := reader.ReadBits(64)
		bitsRead += 64
		if err != nil {
			return errors.New("[ERROR] Couldn't read file length")
		}
		// initialize current_chunk
		var currentChunk bytes.Buffer = *bytes.NewBuffer(make([]byte, 0))
		currentChunkWriter := bitstream.NewWriter(&currentChunk)
		currentChunkLen := 0
		var decompressedBuffer bytes.Buffer = *bytes.NewBuffer(make([]byte, 0))
		decompressedBufferWriter := bitstream.NewWriter(&decompressedBuffer)
		for j := 0; j < int(bufferLen); j++ {
			// Get a bit from the compressed buffer
			bit, err := reader.ReadBit()
			bitsRead++
			if err != nil {
				return errors.New("[ERROR] Couldn't seek past file buffer")
			}
			// Write bit to current chunk
			err = currentChunkWriter.WriteBit(bit)
			if err != nil {
				return errors.New("[ERROR] Failed to write bit to stream")
			}
			// Increment length
			currentChunkLen++
			// Fill the rest with 0s
			err = currentChunkWriter.Flush(bitstream.Zero)
			if err != nil {
				return errors.New("[ERROR] Failed to flush bitstream")
			}
			// Perform htree lookup
			data, foundLeaf, err := decompressor.tree.Lookup(currentChunk, currentChunkLen)
			if err != nil {
				fmt.Println(err)
				return errors.New("[ERROR] Overflow occurred")
			}
			if foundLeaf {
				err := decompressedBufferWriter.WriteByte(data)
				if err != nil {
					return errors.New("[ERROR] Failed to write byte to decompressed buffer")
				}
				currentChunk = *bytes.NewBuffer(make([]byte, 0))
				currentChunkWriter = bitstream.NewWriter(&currentChunk)
			} else {
				// Create a new buffer
				var newBuffer bytes.Buffer = *bytes.NewBuffer(make([]byte, 0))
				newWriter := bitstream.NewWriter(&newBuffer)
				currentChunkReader := bitstream.NewReader(&currentChunk)
				// Copy current_chunk to new_buffer with no padding
				for k := 0; k < currentChunkLen; k++ {
					bit, err = currentChunkReader.ReadBit()
					if err != nil {
						return errors.New("[ERROR] Couldn't transfer data")
					}
					err := newWriter.WriteBit(bit)
					if err != nil {
						return errors.New("[ERROR] Failed to write bit to stream")
					}
				}
				currentChunk = newBuffer
				currentChunkWriter = bitstream.NewWriter(&currentChunk)
			}
		}
		// Create and write file
		dirPath := filepath.Dir(filenameBuffer.String()) // split here
		err = os.MkdirAll(dirPath, 0o755)                // TODO track modes in archive
		if err != nil {
			return errors.New("[ERROR] Couldn't create directory " + dirPath)
		}
		file, err := os.Create(filenameBuffer.String())
		if err != nil {
			return errors.New("[ERROR] Couldn't open file " + filenameBuffer.String())
		}
		_, err = file.Write(decompressedBuffer.Bytes())
		if err != nil {
			return errors.New("[ERROR] Failed to write to file")
		}
		err = file.Close()
		if err != nil {
			return errors.New("[ERROR] Failed to close file")
		}

		// Reset byte boundary
		if bitsRead%8 != 0 {
			bits, err := reader.ReadBits(8 - (bitsRead % 8))
			if bits > 0 {
				return errors.New("[ERROR] Expected bits to be zero")
			}
			if err != nil {
				return errors.New("[ERROR] Couldn't zero out buffer")
			}
		}
	}
	err = bar.Finish()
	if err != nil {
		return errors.New("[ERROR] Failed to cleanly finish progress bar")
	}
	return nil
}
