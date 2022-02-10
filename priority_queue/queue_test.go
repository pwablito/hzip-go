package priority_queue

import (
	"testing"
)

type QueueItemImpl struct {
	Priority int
}

func (hqi QueueItemImpl) Less(item QueueItem) bool {
	return hqi.Priority < item.(QueueItemImpl).Priority
}

func TestHtreePriorityQueue(t *testing.T) {
	item_1 := QueueItemImpl{
		Priority: 1,
	}
	item_3 := QueueItemImpl{
		Priority: 3,
	}
	item_5 := QueueItemImpl{
		Priority: 5,
	}
	item_8 := QueueItemImpl{
		Priority: 8,
	}
	pq := NewPriorityQueue()
	pq.Push(item_5)
	pq.Push(item_8)
	pq.Push(item_3)

	first := pq.Front()
	if first != item_3 {
		t.Error("first should be 3")
		return
	}
	first = pq.Pop()
	if first != item_3 {
		t.Error("first should be 3")
		return
	}
	second := pq.Pop()
	if second != item_5 {
		t.Error("second should be 5")
		return
	}
	pq.Push(item_1)
	length := pq.Length()
	if length != 2 {
		t.Error("length should be 2")
		return
	}
	third := pq.Front()
	if third != item_1 {
		t.Error("third should be 1")
		return
	}
	third = pq.Pop()
	if third != item_1 {
		t.Error("third should be 1")
		return
	}
	fourth := pq.Pop()
	if fourth != item_8 {
		t.Error("fourth should be 8")
		return
	}
	length = pq.Length()
	if length != 0 {
		t.Error("empty length should be 0")
		return
	}
}
