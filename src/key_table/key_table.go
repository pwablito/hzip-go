package key_table

import (
	"bytes"
	"errors"
	"hzip/src/huffman_tree"

	"github.com/dgryski/go-bitstream"
)

type KeyTableData struct {
	length int // bits
	data   string
}

type KeyTable struct {
	table map[string]KeyTableData
}

func (table *KeyTable) Add(key bytes.Buffer, data string, length int) {
	table.table[key.String()] = KeyTableData{
		length: length,
		data:   data,
	}
}

func (table KeyTable) Get(key bytes.Buffer, length int) (*KeyTableData, error) {
	item, ok := table.table[key.String()]
	if !ok {
		return nil, errors.New("couldn't find item")
	}
	return &item, nil
}

func (table *KeyTable) ReadTree(tree *huffman_tree.HuffmanTree) error {
	var buf bytes.Buffer
	table.AddSubtreeWithPrefix(buf, 0, &tree.Head)
	return nil
}

func (table *KeyTable) AddSubtreeWithPrefix(prefix bytes.Buffer, prefix_len int, tree_node *huffman_tree.HTreeNode) {
	if (*tree_node).IsLeaf() {
		table.Add(prefix, string((*tree_node).Data()), prefix_len)
	} else {
		if (*tree_node).Left() != nil {
			var buf bytes.Buffer
			writer := bitstream.NewWriter(&buf)
			CopyNumBitsToBitstreamWriter(prefix, writer, prefix_len)
			writer.WriteBit(bitstream.Zero)
			writer.Flush(bitstream.Zero)
			table.AddSubtreeWithPrefix(buf, prefix_len+1, (*tree_node).Left())
		}
		if (*tree_node).Right() != nil {
			var buf bytes.Buffer
			writer := bitstream.NewWriter(&buf)
			CopyNumBitsToBitstreamWriter(prefix, writer, prefix_len)
			writer.WriteBit(bitstream.One)
			writer.Flush(bitstream.Zero)
			table.AddSubtreeWithPrefix(buf, prefix_len+1, (*tree_node).Left())
		}
	}
}

func CopyNumBitsToBitstreamWriter(in bytes.Buffer, out *bitstream.BitWriter, num_bits int) error {
	in_stream := bitstream.NewReader(bytes.NewReader(in.Bytes()))
	for i := 0; i < num_bits; i++ {
		bit, err := in_stream.ReadBit()
		if err != nil {
			return err
		}
		out.WriteBit(bit)
	}
	return nil
}
