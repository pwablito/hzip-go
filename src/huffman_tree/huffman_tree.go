package huffman_tree

import (
	"bytes"
	"errors"

	"github.com/dgryski/go-bitstream"
)

type HuffmanTree struct {
	Head      HTreeNode
	Frequency int
}

func (tree *HuffmanTree) Lookup(buf bytes.Buffer, length int) (byte, bool, error) {
	/*
		If found, return (data, true, nil)
		If not overflow but not target leaf, return (0, false, nil)
		If overflow return (0, false, error)
	*/
	return _lookup(tree.Head, &buf, length)
}

func _lookup(node HTreeNode, buf *bytes.Buffer, length int) (byte, bool, error) {
	if length == 0 {
		if !node.IsLeaf() {
			return 0, false, nil
		}
		return node.Data(), true, nil
	}
	if node.IsLeaf() {
		return 0, false, errors.New("[ERROR] Tree overflow")
	}
	node_as_tree := node.(TreeNode)
	buf_reader := bitstream.NewReader(buf)
	bit, err := buf_reader.ReadBit()
	var next_buffer bytes.Buffer
	next_buffer_writer := bitstream.NewWriter(&next_buffer)
	for i := 0; i < length-1; i++ {
		next_bit, _ := buf_reader.ReadBit()
		next_buffer_writer.WriteBit(next_bit)
	}
	next_buffer_writer.Flush(bitstream.Zero)
	if err != nil {
		return 0, false, errors.New("[ERROR] Buffer overflow")
	}
	if bit == bitstream.One {
		return _lookup(*(node_as_tree.Right()), &next_buffer, length-1)
	} else if bit == bitstream.Zero {
		return _lookup(*(node_as_tree.Left()), &next_buffer, length-1)
	} else {
		return 0, false, errors.New("[ERROR] Invalid bit")
	}
}

type HTreeNode interface {
	IsLeaf() bool
	Frequency() int
	Data() byte
	Left() *HTreeNode
	Right() *HTreeNode
}

func CombineTrees(tree_1 *HuffmanTree, tree_2 *HuffmanTree) *HuffmanTree {
	new_tree := HuffmanTree{
		Frequency: tree_1.Frequency + tree_2.Frequency,
	}
	head_node := TreeNode{}
	if tree_1.Frequency < tree_2.Frequency {
		head_node.LeftChild = &tree_1.Head
		head_node.RightChild = &tree_2.Head
	} else {
		head_node.LeftChild = &tree_2.Head
		head_node.RightChild = &tree_1.Head
	}
	new_tree.Head = head_node
	return &new_tree
}

type LeafNode struct {
	Freq     int
	LeafData byte
}

func (node LeafNode) IsLeaf() bool {
	return true
}

func (node LeafNode) Frequency() int {
	return node.Freq
}

func (node LeafNode) Data() byte {
	return node.LeafData
}

func (node LeafNode) Left() *HTreeNode {
	return nil
}
func (node LeafNode) Right() *HTreeNode {
	return nil
}

type TreeNode struct {
	LeftChild  *HTreeNode
	RightChild *HTreeNode
}

func (node TreeNode) IsLeaf() bool {
	return false
}

func (node TreeNode) Frequency() int {
	return -1
}

func (node TreeNode) Data() byte {
	return 0
}

func (node TreeNode) Left() *HTreeNode {
	return node.LeftChild
}
func (node TreeNode) Right() *HTreeNode {
	return node.RightChild
}
