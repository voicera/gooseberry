// Package queues provides queue data structures.
package queues

import (
	"container/heap"
)

// PriorityQueue represents a priority queue. This type is NOT thread-safe!
type PriorityQueue interface {
	// Add adds the specified item to the priority queue.
	// Priority is a integer that represents the specified item's relative
	// priority: the bigger, the higher.
	Add(item interface{}, priority int)

	// Poll removes the lowest-priority item from the queue and returns it.
	Poll() (interface{}, bool)

	// Len is the number of items in the queue.
	Len() int
}

type prioritizedItem struct {
	item     interface{}
	priority int
	index    int
}

type priorityQueue []*prioritizedItem

// NewPriorityQueue creates a new instance of PriorityQueue.
func NewPriorityQueue(initialCapacity int) PriorityQueue {
	pq := make(priorityQueue, 0, initialCapacity)
	heap.Init(&pq)
	return &pq
}

func (pq *priorityQueue) Add(item interface{}, priority int) {
	heap.Push(pq, &prioritizedItem{item: item, priority: priority})
}

func (pq *priorityQueue) Push(prioritized interface{}) {
	item := prioritized.(*prioritizedItem)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Poll() (interface{}, bool) {
	if len(*pq) == 0 {
		return nil, false
	}
	return heap.Pop(pq).(*prioritizedItem).item, true
}

func (pq *priorityQueue) Pop() interface{} {
	indexOfLastItem := len(*pq) - 1
	item := (*pq)[indexOfLastItem]
	item.index = -1 // for safety
	*pq = (*pq)[:indexOfLastItem]
	return item
}

func (pq *priorityQueue) Len() int {
	return len(*pq)
}

func (pq *priorityQueue) Less(i int, j int) bool {
	return (*pq)[i].priority > (*pq)[j].priority
}

func (pq *priorityQueue) Swap(i int, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].index = i
	(*pq)[j].index = j
}
