package compression

import "hzip/util"

type HtreeQueueItem struct {
	Priority int
	Tree     *HuffmanTree
}

func (hqi HtreeQueueItem) Less(item util.QueueItem) bool {
	return hqi.Priority < item.(HtreeQueueItem).Priority
}
