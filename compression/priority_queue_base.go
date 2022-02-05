// NOTE: Taken from the internet

package compression

import "container/heap"

// QueueItem represents the interface to be implemented stored in this queue
type QueueItem interface {
	Less(item QueueItem) bool
}

// priorityQueueImpl for the underlying implementation of priority queues
type priorityQueueImpl []QueueItem

// Len get queue length
func (pqi priorityQueueImpl) Len() int {
	return len(pqi)
}

// Less is used for element comparison
func (pqi priorityQueueImpl) Less(i, j int) bool {
	return pqi[i].Less(pqi[j])
}

// Swap
func (pqi priorityQueueImpl) Swap(i, j int) {
	pqi[i], pqi[j] = pqi[j], pqi[i]
}

// Push is used to push an object into the queue
func (pqi *priorityQueueImpl) Push(x interface{}) {
	item := x.(QueueItem)
	*pqi = append(*pqi, item)
}

// Pop pops an object out of the queue
func (pqi *priorityQueueImpl) Pop() interface{} {
	old := *pqi
	n := len(old)
	item := old[n-1]
	*pqi = old[0 : n-1]
	return item
}

// PriorityQueue implements priority queue
type PriorityQueue struct {
	priorityQueueImpl
}

// NewPriorityQueue is used to build PriorityQueue
func NewPriorityQueue() *PriorityQueue {
	var pq PriorityQueue
	heap.Init(&pq.priorityQueueImpl)
	return &pq
}

// Push is used to push an object into the queue
func (pq *PriorityQueue) Push(item QueueItem) {
	heap.Push(&pq.priorityQueueImpl, item)
}

// Pop is used to pop an object from the queue
func (pq *PriorityQueue) Pop() QueueItem {
	return heap.Pop(&pq.priorityQueueImpl).(QueueItem)
}

// Front is used to get the minimum value in the current queue
func (pq *PriorityQueue) Front() QueueItem {
	// The first bit in the queue should be the minimum
	return pq.priorityQueueImpl[0]
}

// Length is used to get the length of the current queue
func (pq *PriorityQueue) Length() int {
	return pq.priorityQueueImpl.Len()
}
