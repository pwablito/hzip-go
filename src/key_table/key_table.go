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

func (table *KeyTable) WriteTree() (*huffman_tree.HuffmanTree, error) {
	tree_head := huffman_tree.TreeNode{
		LeftChild:  nil,
		RightChild: nil,
	}
	for key, value := range table.Table {
		var current_node *huffman_tree.TreeNode
		current_node = &tree_head
		reader := bitstream.NewReader(&value.Data)
		for i := 0; i < value.Length-1; i++ {
			bit, err := reader.ReadBit()
			if err != nil {
				return nil, errors.New("[ERROR] Failed to read bit from key table value")
			}
			var blank_tree_node huffman_tree.HTreeNode = huffman_tree.TreeNode{
				LeftChild:  nil,
				RightChild: nil,
			}
			if bit == bitstream.Zero { // Left
				if current_node.LeftChild == nil {
					current_node.LeftChild = &blank_tree_node
				}
				current_node = (*current_node.LeftChild).(*huffman_tree.TreeNode)
			} else if bit == bitstream.One { // Right
				if current_node.LeftChild == nil {
					current_node.RightChild = &blank_tree_node
				}
				current_node = (*current_node.RightChild).(*huffman_tree.TreeNode)
			} else {
				return nil, errors.New("[ERROR] Got invalid bit from key table value")
			}
		}
		// current_node is now where we will append a leaf node
		bit, err := reader.ReadBit()
		if err != nil {
			return nil, errors.New("[ERROR] Failed to read bit from key table value")
		}
		var new_leaf huffman_tree.HTreeNode = huffman_tree.LeafNode{
			Freq:     0, // We don't care about frequency at this point
			LeafData: key,
		}
		if bit == bitstream.Zero { // Left
			current_node.LeftChild = &new_leaf
		} else if bit == bitstream.One { // Right
			current_node.LeftChild = &new_leaf
		} else {
			return nil, errors.New("[ERROR] Got invalid bit from key table value")
		}

	}
	tree := huffman_tree.HuffmanTree{
		Head:      tree_head,
		Frequency: 0, // We don't care about frequency when we are writing the tree
	}
	return &tree, nil
}

func (table *KeyTable) AddSubtreeWithPrefix(prefix bytes.Buffer, prefix_len int, tree_node *huffman_tree.HTreeNode) {
	if (*tree_node).IsLeaf() {
		table.Add((*tree_node).Data(), prefix, prefix_len)
	} else {
		if (*tree_node).Left() != nil {
			// Left is a 0
			var buf bytes.Buffer
			writer := bitstream.NewWriter(&buf)
			CopyNumBitsToBitstreamWriter(prefix, writer, prefix_len)
			writer.WriteBit(bitstream.Zero)
			writer.Flush(bitstream.Zero)
			table.AddSubtreeWithPrefix(buf, prefix_len+1, (*tree_node).Left())
		}
		if (*tree_node).Right() != nil {
			// Right is a 1
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
