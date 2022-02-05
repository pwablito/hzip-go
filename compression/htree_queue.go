package compression

type HtreeQueueItem struct {
	Priority int
	// Tree     HuffmanTree
}

func (hqi HtreeQueueItem) Less(item QueueItem) bool {
	return hqi.Priority < item.(HtreeQueueItem).Priority
}
