package compression

type HuffmanTree struct {
	Head      HTreeNode
	Frequency int
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
