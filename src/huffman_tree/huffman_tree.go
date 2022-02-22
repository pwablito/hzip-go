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
	reader := bitstream.NewReader(&buf)
	var current_node HTreeNode = tree.Head
	for i := 0; i < length; i++ {
		if current_node.IsLeaf() {
			return 0, false, errors.New("[ERROR] Tree overflow")
		}
		bit, err := reader.ReadBit()
		if err != nil {
			return 0, false, errors.New("[ERROR] Buffer overflow")
		}
		if bit == bitstream.One { // right
			current_node = *current_node.(TreeNode).RightChild
		} else if bit == bitstream.Zero { // left
			current_node = *current_node.(TreeNode).LeftChild
		} else {
			return 0, false, errors.New("[ERROR] Invalid bit")
		}
	}
	if !current_node.IsLeaf() {
		return 0, false, nil
	}
	return current_node.Data(), true, nil
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
