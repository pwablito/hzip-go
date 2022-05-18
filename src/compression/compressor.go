package compression

import (
	"bytes"
	"errors"
	"fmt"
	"hzip/src/frequency_table"
	"hzip/src/huffman_tree"
	"hzip/src/input"
	"hzip/src/key_table"
	"hzip/src/output"
	"hzip/src/priority_queue"
	"os"

	"github.com/dgryski/go-bitstream"
	"github.com/schollz/progressbar/v3"
)

type Compressor struct {
	Output   output.Output
	Inputs   []input.Input
	keyTable key_table.KeyTable
}

func (compressor *Compressor) GenerateScheme() error {
	fmt.Println("[INFO] Creating frequency table")
	freqTable := frequency_table.CreateFrequencyTable()
	bar := progressbar.NewOptions(
		len(compressor.Inputs),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetPredictTime(true),
	)
	for _, inputObj := range compressor.Inputs {
		err := bar.Add(1)
		if err != nil {
			return errors.New("[ERROR] Failed to update progress bar")
		}
		data, err := inputObj.GetData()
		if err != nil {
			fmt.Println(err)
			return errors.New("[ERROR] Failed to read data from input")
		}
		for _, currentByte := range data {
			freqTable.Increment(currentByte)
		}
	}
	err := bar.Finish()
	if err != nil {
		return errors.New("[ERROR] Failed to finish progress bar")
	}
	fmt.Println("[INFO] Constructing Huffman Tree")
	pq := priority_queue.NewPriorityQueue()
	for data, frequency := range freqTable.GetFrequencies() {
		pq.Push(huffman_tree.HtreeQueueItem{
			Priority: frequency,
			Tree: &huffman_tree.HuffmanTree{
				Head: huffman_tree.LeafNode{
					Freq:     frequency,
					LeafData: data,
				},
				Frequency: frequency,
			},
		})
	}
	for pq.Len() > 1 {
		newTree := huffman_tree.CombineTrees(pq.Pop().(huffman_tree.HtreeQueueItem).Tree, pq.Pop().(huffman_tree.HtreeQueueItem).Tree)
		pq.Push(huffman_tree.HtreeQueueItem{
			Priority: newTree.Frequency,
			Tree:     newTree,
		})
	}
	finalTree := pq.Pop().(huffman_tree.HtreeQueueItem).Tree
	compressor.keyTable = key_table.CreateKeyTable()
	err = compressor.keyTable.ReadTree(finalTree)
	if err != nil {
		fmt.Println(err)
		return errors.New("[ERROR] Failed to generate keys from Huffman tree")
	}
	return nil
}

func (compressor *Compressor) CompressToOutput() error {
	/*
		Output looks like this:
		----------------------------------------------
		|--- Number of key table entries (64 bits) ---|
		for each key table entry {
			|--- key (1 byte) ---|
			|--- length (8 bytes) ---|
			|--- value ($length bits) ---|
		}

		|--- 0 until edge of byte boundary ---|

		|--- number of inputs (8 bytes) ---|
		for each input {
			|--- length of filename (8 bytes) ---|
			|--- filename ($length bytes) ---|
			|--- length of compressed buffer (8 bytes)---|
			|--- compressed buffer ($length bits) ---|
			|--- 0 until edge of byte boundary ---|
		}
		----------------------------------------------
	*/
	err := compressor.Output.Open()
	if err != nil {
		fmt.Println(err)
		return errors.New("[ERROR] Couldn't open output")
	}
	defer func(Output output.Output) {
		err := Output.Close()
		if err != nil {
			fmt.Println("[FATAL] Failed to close output")
			os.Exit(1)
		}
	}(compressor.Output)
	bar := progressbar.NewOptions(
		len(compressor.Inputs),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetPredictTime(true),
	)
	// Dump key table to output
	var keyTableBuffer bytes.Buffer
	keyTableWriter := bitstream.NewWriter(&keyTableBuffer)
	err = keyTableWriter.WriteBits(uint64(len(compressor.keyTable.Table)), 64)
	if err != nil {
		return errors.New("[ERROR] Failed to write bits to stream")
	}
	for key, value := range compressor.keyTable.Table {
		err := keyTableWriter.WriteByte(key)
		if err != nil {
			return errors.New("[ERROR] Failed to write byte to key table")
		}
		err = keyTableWriter.WriteBits(uint64(value.Length), 64)
		if err != nil {
			return errors.New("[ERROR] Failed to write bits to key table")
		}
		reader := bitstream.NewReader(&value.Data)
		for i := 0; i < value.Length; i++ {
			bit, err := reader.ReadBit()
			if err != nil {
				return errors.New("[ERROR] Failed to read bit from key table entry")
			}
			err = keyTableWriter.WriteBit(bit)
			if err != nil {
				return errors.New("[ERROR] Failed to write bit to key table")
			}
		}
	}
	err = keyTableWriter.Flush(bitstream.Zero)
	if err != nil {
		return errors.New("[ERROR] Failed to flush bitstream")
	}
	err = compressor.Output.Write(keyTableBuffer.Bytes())
	if err != nil {
		return errors.New("[ERROR] Failed to write to output buffer")
	}
	var numInputsBuffer bytes.Buffer
	numInputsWriter := bitstream.NewWriter(&numInputsBuffer)
	// TODO This doesn't need to be a 64 bit int
	err = numInputsWriter.WriteBits(uint64(len(compressor.Inputs)), 64)
	if err != nil {
		return errors.New("[ERROR] Failed to write bits to number input writer")
	}
	err = compressor.Output.Write(numInputsBuffer.Bytes())
	if err != nil {
		return errors.New("[ERROR] Failed to write bytes to compressor output")
	}
	for _, inputObj := range compressor.Inputs {
		err := bar.Add(1)
		if err != nil {
			return errors.New("[ERROR] Failed to update progress bar")
		}
		inputData, err := inputObj.GetData()
		if err != nil {
			fmt.Println(err)
			return errors.New("[ERROR] Failed to get data from input")
		}
		contentBuffer, compressedBufferLen, err := compressor.compress_buffer(inputData)
		if err != nil {
			fmt.Println(err)
			return errors.New("[ERROR] Failed to compress buffer")
		}
		var metaBuffer bytes.Buffer
		// TODO Compress filenames too
		metaWriter := bitstream.NewWriter(&metaBuffer)
		err = metaWriter.WriteBits(uint64(len(inputObj.(input.FileInput).Filename)), 64)
		if err != nil {
			return errors.New("[ERROR] Failed to write metadata bits to buffer")
		}
		for _, character := range inputObj.(input.FileInput).Filename {
			err := metaWriter.WriteByte(byte(character))
			if err != nil {
				return errors.New("[ERROR] Failed to write metadata bytes")
			}
		}
		err = metaWriter.WriteBits(uint64(compressedBufferLen), 64)
		if err != nil {
			return errors.New("[ERROR] Failed to write bits to metadata buffer")
		}
		err = compressor.Output.Write(metaBuffer.Bytes())
		if err != nil {
			fmt.Println(err)
			return errors.New("[ERROR] Failed to write metadata to output")
		}
		err = compressor.Output.Write(contentBuffer.Bytes())
		if err != nil {
			fmt.Println(err)
			return errors.New("[ERROR] Failed to write compressed buffer to output")
		}
	}
	err = bar.Finish()
	if err != nil {
		return errors.New("[ERROR] Failed to finish progress bar")
	}
	return nil
}

func (compressor *Compressor) compress_buffer(input_buffer []byte) (*bytes.Buffer, int, error) {
	var outputBuffer bytes.Buffer
	totalBits := 0
	outputWriter := bitstream.NewWriter(&outputBuffer)
	for _, currentByte := range input_buffer {
		data, err := compressor.keyTable.Get(currentByte)
		if err != nil {
			return nil, 0, err
		}
		reader := bitstream.NewReader(&data.Data)
		for i := 0; i < data.Length; i++ {
			nextBit, err := reader.ReadBit()
			if err != nil {
				fmt.Println(err)
				return nil, 0, errors.New("[ERROR] Couldn't read data from reader bitstream")
			}
			err = outputWriter.WriteBit(nextBit)
			totalBits++
			if err != nil {
				fmt.Println(err)
				return nil, 0, errors.New("[ERROR] Couldn't write data to writer bitstream")
			}
		}
	}
	err := outputWriter.Flush(bitstream.Zero)
	if err != nil {
		return nil, 0, errors.New("[ERROR] Failed to flush bitstream")
	}
	return &outputBuffer, totalBits, nil
}

func (compressor *Compressor) AddInput(inputObj input.Input) {
	compressor.Inputs = append(compressor.Inputs, inputObj)
}

func (compressor *Compressor) SetOutput(outputObj output.Output) {
	compressor.Output = outputObj
}
