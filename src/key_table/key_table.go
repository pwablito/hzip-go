package key_table

import (
	"bytes"
	"errors"
	"hzip/src/huffman_tree"

	"github.com/dgryski/go-bitstream"
)

type KeyTableData struct {
	Length int // bits
	Data   bytes.Buffer
}

type KeyTable struct {
	Table map[byte]KeyTableData
}

func (table *KeyTable) Add(key byte, data bytes.Buffer, length int) {
	table.Table[key] = KeyTableData{
		Length: length,
		Data:   data,
	}
}

func (table KeyTable) Get(key byte) (*KeyTableData, error) {
	item, ok := table.Table[key]
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
		table.Add((*tree_node).Data(), prefix, prefix_len)
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
			table.AddSubtreeWithPrefix(buf, prefix_len+1, (*tree_node).Right())
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
