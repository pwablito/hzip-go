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
	nodeAsTree := node.(TreeNode)
	bufReader := bitstream.NewReader(buf)
	bit, err := bufReader.ReadBit()
	if err != nil {
		return 0, false, errors.New("[ERROR] Failed to read bit from buffer")
	}
	var nextBuffer bytes.Buffer
	nextBufferWriter := bitstream.NewWriter(&nextBuffer)
	for i := 0; i < length-1; i++ {
		nextBit, _ := bufReader.ReadBit()
		err := nextBufferWriter.WriteBit(nextBit)
		if err != nil {
			return 0, false, errors.New("[ERROR] Failed to write to temporary buffer")
		}
	}
	err = nextBufferWriter.Flush(bitstream.Zero)
	if err != nil {
		return 0, false, errors.New("[ERROR] Failed to flush buffer")
	}
	if err != nil {
		return 0, false, errors.New("[ERROR] Buffer overflow")
	}
	if bit == bitstream.One {
		return _lookup(*(nodeAsTree.Right()), &nextBuffer, length-1)
	} else if bit == bitstream.Zero {
		return _lookup(*(nodeAsTree.Left()), &nextBuffer, length-1)
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

func CombineTrees(tree1 *HuffmanTree, tree2 *HuffmanTree) *HuffmanTree {
	newTree := HuffmanTree{
		Frequency: tree1.Frequency + tree2.Frequency,
	}
	headNode := TreeNode{}
	if tree1.Frequency < tree2.Frequency {
		headNode.LeftChild = &tree1.Head
		headNode.RightChild = &tree2.Head
	} else {
		headNode.LeftChild = &tree2.Head
		headNode.RightChild = &tree1.Head
	}
	newTree.Head = headNode
	return &newTree
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
