package queue

import (
	"container/heap"
	"github.com/yates-z/easel/core/pool"
	"sync"
)

// Item defines the elements in the priority queue
type Item[T any] struct {
	value    T   // The value of the item
	priority int // The priority of the item
}

// PriorityQueue defines the priority queue
type PriorityQueue[T any] struct {
	items    []*Item[T]
	mu       sync.RWMutex
	lessFunc func(a, b *Item[T]) bool
	itemPool *pool.Pool[*Item[T]]
}

func NewPriorityQueue[T any]() *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		lessFunc: func(a, b *Item[T]) bool { return a.priority < b.priority },
		itemPool: pool.New(func() *Item[T] {
			return new(Item[T])
		}),
	}
	return pq
}

// Len implements the Len method of sort.Interface
func (pq *PriorityQueue[T]) Len() int { return len(pq.items) }

// Less implements the Less method of sort.Interface
// Items with higher priority will appear earlier
func (pq *PriorityQueue[T]) Less(i, j int) bool {
	return pq.lessFunc(pq.items[i], pq.items[j])
}

// Swap implements the Swap method of sort.Interface
func (pq *PriorityQueue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

// Push adds an element to the heap
func (pq *PriorityQueue[T]) Push(x any) {
	item := x.(*Item[T])
	pq.items = append(pq.items, item)
}

// Pop removes and returns the highest-priority element from the heap
func (pq *PriorityQueue[T]) Pop() any {
	old := pq.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	pq.items = old[0 : n-1]
	return item
}

func (pq *PriorityQueue[T]) Enqueue(value T, priority int) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	item := pq.itemPool.Get()
	item.value = value
	item.priority = priority
	heap.Push(pq, item)
}

func (pq *PriorityQueue[T]) Dequeue() (T, bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if pq.Len() == 0 {
		var zeroValue T
		return zeroValue, false
	}
	item := heap.Pop(pq).(*Item[T])
	value := item.value
	item.priority = 0
	item.value = *new(T)
	pq.itemPool.Put(item)
	return value, true
}

func (pq *PriorityQueue[T]) Peek() (T, bool) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	if pq.Len() == 0 {
		var zeroValue T
		return zeroValue, false
	}
	return pq.items[0].value, true
}
