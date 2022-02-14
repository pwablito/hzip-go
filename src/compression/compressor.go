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

	"github.com/dgryski/go-bitstream"
	"github.com/schollz/progressbar/v3"
)

type Compressor struct {
	Output          output.Output
	Inputs          []input.Input
	frequency_table frequency_table.FrequencyTable
	key_table       key_table.KeyTable
}

func (compressor *Compressor) GenerateScheme() error {
	fmt.Println("[INFO] Creating frequency table")
	bar := progressbar.NewOptions(
		len(compressor.Inputs),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetPredictTime(true),
	)
	for _, input_obj := range compressor.Inputs {
		bar.Add(1)
		data, err := input_obj.GetData()
		if err != nil {
			fmt.Println(err)
			return errors.New("[ERROR] Failed to read data from input")
		}
		for _, current_byte := range data {
			compressor.frequency_table.Increment(current_byte)
		}
	}
	bar.Finish()
	fmt.Println("[INFO] Constructing Huffman Tree")
	pq := priority_queue.NewPriorityQueue()
	for data, frequency := range compressor.frequency_table.GetFrequencies() {
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
		new_tree := huffman_tree.CombineTrees(pq.Pop().(huffman_tree.HtreeQueueItem).Tree, pq.Pop().(huffman_tree.HtreeQueueItem).Tree)
		pq.Push(huffman_tree.HtreeQueueItem{
			Priority: new_tree.Frequency,
			Tree:     new_tree,
		})
	}
	final_tree := pq.Pop().(huffman_tree.HtreeQueueItem).Tree
	compressor.key_table = key_table.CreateKeyTable()
	err := compressor.key_table.ReadTree(final_tree)
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
	defer compressor.Output.Close()
	bar := progressbar.NewOptions(
		len(compressor.Inputs),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetPredictTime(true),
	)
	// Dump key table to output
	var key_table_buffer bytes.Buffer
	key_table_writer := bitstream.NewWriter(&key_table_buffer)
	key_table_writer.WriteBits(uint64(len(compressor.key_table.Table)), 64)
	for key, value := range compressor.key_table.Table {
		key_table_writer.WriteByte(key)
		key_table_writer.WriteBits(uint64(value.Length), 64)
		reader := bitstream.NewReader(&value.Data)
		for i := 0; i < value.Length; i++ {
			bit, err := reader.ReadBit()
			if err != nil {
				return errors.New("[ERROR] Failed to read bit from key table entry")
			}
			key_table_writer.WriteBit(bit)
		}
	}
	key_table_writer.Flush(bitstream.Zero)
	compressor.Output.Write(key_table_buffer.Bytes())
	var num_inputs_buffer bytes.Buffer
	num_inputs_writer := bitstream.NewWriter(&num_inputs_buffer)
	num_inputs_writer.WriteBits(uint64(len(compressor.Inputs)), 64)
	compressor.Output.Write(num_inputs_buffer.Bytes())
	for _, input_obj := range compressor.Inputs {
		bar.Add(1)
		input_data, err := input_obj.GetData()
		if err != nil {
			fmt.Println(err)
			return errors.New("[ERROR] Failed to get data from input")
		}
		content_buffer, _, err := compressor.compress_buffer(input_data)
		if err != nil {
			fmt.Println(err)
			return errors.New("[ERROR] Failed to compress buffer")
		}
		var meta_buffer bytes.Buffer
		meta_writer := bitstream.NewWriter(&meta_buffer)
		meta_writer.WriteBits(uint64(len(input_obj.(input.FileInput).Filename)), 64)
		for _, character := range input_obj.(input.FileInput).Filename {
			meta_writer.WriteByte(byte(character))
		}
		meta_writer.WriteBits(uint64(content_buffer.Len()), 64)
		err = compressor.Output.Write(meta_buffer.Bytes())
		if err != nil {
			fmt.Println(err)
			return errors.New("[ERROR] Failed to write metadata to output")
		}
		err = compressor.Output.Write(content_buffer.Bytes())
		if err != nil {
			fmt.Println(err)
			return errors.New("[ERROR] Failed to write compressed buffer to output")
		}
	}
	bar.Finish()
	return nil
}

func (compressor *Compressor) compress_buffer(input_buffer []byte) (*bytes.Buffer, int, error) {
	var output_buffer bytes.Buffer
	total_bits := 0
	output_writer := bitstream.NewWriter(&output_buffer)
	for _, current_byte := range input_buffer {
		data, err := compressor.key_table.Get(current_byte)
		if err != nil {
			return nil, 0, err
		}
		reader := bitstream.NewReader(&data.Data)
		for i := 0; i < data.Length; i++ {
			next_bit, err := reader.ReadBit()
			if err != nil {
				fmt.Println(err)
				return nil, 0, errors.New("[ERROR] Couldn't read data from reader bitstream")
			}
			err = output_writer.WriteBit(next_bit)
			if err != nil {
				fmt.Println(err)
				return nil, 0, errors.New("[ERROR] Couldn't write data to writer bitstream")
			}
		}
	}
	output_writer.Flush(bitstream.Zero)
	return &output_buffer, total_bits, nil
}

func (compressor *Compressor) AddInput(input_obj input.Input) {
	compressor.Inputs = append(compressor.Inputs, input_obj)
}

func (compressor *Compressor) SetOutput(output_obj output.Output) {
	compressor.Output = output_obj
}
