package huffman_tree

import "hzip/priority_queue"

type HtreeQueueItem struct {
	Priority int
	Tree     *HuffmanTree
}

func (hqi HtreeQueueItem) Less(item priority_queue.QueueItem) bool {
	return hqi.Priority < item.(HtreeQueueItem).Priority
}
