package compression

type KeyTableData struct {
	length int // bits
	data   string
}

type KeyTable struct {
	table map[string]KeyTableData
}

func (table *KeyTable) Add(key string, data string, length int) {
	table.table[key] = KeyTableData{
		length: length,
		data:   data,
	}
}

func (table *KeyTable) ReadTree(tree *HuffmanTree) error {
	table.AddSubtreeWithPrefix("", &tree.Head)
	return nil
}

func (table *KeyTable) AddSubtreeWithPrefix(prefix string, tree_node *HTreeNode) {
	if (*tree_node).IsLeaf() {
		table.Add(prefix, string((*tree_node).Data()), len(prefix))
	} else {
		if (*tree_node).Left() != nil {
			table.AddSubtreeWithPrefix(prefix+"0", (*tree_node).Left())
		}
		if (*tree_node).Right() != nil {
			table.AddSubtreeWithPrefix(prefix+"1", (*tree_node).Right())
		}
	}
}
