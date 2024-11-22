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

type Option[T any] func(q *PriorityQueue[T])

func WithLessFunc[T any](less func(a, b *Item[T]) bool) Option[T] {
	return func(q *PriorityQueue[T]) {
		q.lessFunc = less
	}
}

// PriorityQueue defines the priority queue
type PriorityQueue[T any] struct {
	items    []*Item[T]
	mu       sync.RWMutex
	lessFunc func(a, b *Item[T]) bool
	itemPool *pool.Pool[*Item[T]]
}

func NewPriorityQueue[T any](opts ...Option[T]) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		lessFunc: func(a, b *Item[T]) bool { return a.priority > b.priority },
		itemPool: pool.New(func() *Item[T] {
			return new(Item[T])
		}),
	}
	for _, opt := range opts {
		opt(pq)
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

// Enqueue method adds elements to the queue.
func (pq *PriorityQueue[T]) Enqueue(value T, priority int) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	item := pq.itemPool.Get()
	item.value = value
	item.priority = priority
	heap.Push(pq, item)
}

// Dequeue pops out of the queue and returns an element.
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

// Peek method returns the next element to be Dequeued,
// but it will not actually be Dequeued.
func (pq *PriorityQueue[T]) Peek() (T, bool) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	if pq.Len() == 0 {
		var zeroValue T
		return zeroValue, false
	}
	return pq.items[0].value, true
}

// Clear all elements in the queue.
func (pq *PriorityQueue[T]) Clear() {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	for _, item := range pq.items {
		item.priority = 0
		pq.itemPool.Put(item)
	}
	pq.items = nil
}

// Iterator returns all elements in the queue.
func (pq *PriorityQueue[T]) Iterator() []T {
	var result []T

	pq.mu.RLock()
	defer pq.mu.RUnlock()

	var items []*Item[T]
	for range pq.Len() {
		item := heap.Pop(pq).(*Item[T])
		result = append(result, item.value)
		items = append(items, item)
	}

	for _, item := range items {
		heap.Push(pq, item)
	}

	return result
}
