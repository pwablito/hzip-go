package compression

import (
	"errors"
	"fmt"
	"hzip/src/frequency_table"
	"hzip/src/huffman_tree"
	"hzip/src/input"
	"hzip/src/key_table"
	"hzip/src/output"
	"hzip/src/priority_queue"

	"github.com/schollz/progressbar/v3"
)

type Compressor struct {
	Output          output.Output
	Inputs          []input.Input
	frequency_table frequency_table.FrequencyTable
}

func (compressor *Compressor) Process() error {
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
	key_table := key_table.CreateKeyTable()
	err := key_table.ReadTree(final_tree)
	if err != nil {
		fmt.Println(err)
		return errors.New("[ERROR] Failed to generate keys from Huffman tree")
	}
	return nil
}

func (compressor *Compressor) Dump() error {
	return errors.New("[ERROR] Not implemented")
}

func (compressor *Compressor) AddInput(input_obj input.Input) {
	compressor.Inputs = append(compressor.Inputs, input_obj)
}

func (compressor *Compressor) SetOutput(output_obj output.Output) {
	compressor.Output = output_obj
}
