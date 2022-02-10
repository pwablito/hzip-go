package compression

import (
	"errors"
	"fmt"
	"hzip/input"
	"hzip/output"
	"hzip/util"

	"github.com/schollz/progressbar/v3"
)

type Compressor struct {
	Output          output.Output
	Inputs          []input.Input
	frequency_table FrequencyTable
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
	pq := util.NewPriorityQueue()
	for data, frequency := range compressor.frequency_table.GetFrequencies() {
		pq.Push(HtreeQueueItem{
			Priority: frequency,
			Tree: &HuffmanTree{
				Head: LeafNode{
					Freq:     frequency,
					LeafData: data,
				},
				Frequency: frequency,
			},
		})
	}
	for pq.Len() > 1 {
		new_tree := CombineTrees(pq.Pop().(HtreeQueueItem).Tree, pq.Pop().(HtreeQueueItem).Tree)
		pq.Push(HtreeQueueItem{
			Priority: new_tree.Frequency,
			Tree:     new_tree,
		})
	}
	final_tree := pq.Pop().(HtreeQueueItem).Tree
	key_table := CreateKeyTable()
	err := key_table.ReadTree(final_tree)
	if err != nil {
		fmt.Println(err)
		return errors.New("[ERROR] Failed to generate keys from Huffman tree")
	}

	return errors.New("[ERROR] Compression not fully implemented")
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
