package queue

import (
	"container/heap"
	"github.com/yates-z/easel/core/pool"
	"reflect"
	"sync"
)

// Item defines the elements in the priority queue
type Item[T any] struct {
	value    T   // The value of the item
	priority int // The priority of the item
}

func (i Item[T]) Value() T {
	return i.value
}

type priorityQueue[T any] struct {
	items     []*Item[T]
	lessFunc  func(a, b *Item[T]) bool
	equalFunc func(a, b T) bool
}

// Len implements the Len method of sort.Interface
func (pq *priorityQueue[T]) Len() int { return len(pq.items) }

// Less implements the Less method of sort.Interface
// Items with higher priority will appear earlier
func (pq *priorityQueue[T]) Less(i, j int) bool {
	return pq.lessFunc(pq.items[i], pq.items[j])
}

// Swap implements the Swap method of sort.Interface
func (pq *priorityQueue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

// Push adds an element to the heap
func (pq *priorityQueue[T]) Push(x any) {
	item := x.(*Item[T])
	pq.items = append(pq.items, item)
}

// Pop removes and returns the highest-priority element from the heap
func (pq *priorityQueue[T]) Pop() any {
	old := pq.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	pq.items = old[0 : n-1]
	return item
}

// PriorityQueue defines the priority queue
type PriorityQueue[T any] struct {
	priorityQueue[T]
	mu       sync.RWMutex
	itemPool *pool.Pool[*Item[T]]
}

type Option[T any] func(q *PriorityQueue[T])

func WithLessFunc[T any](less func(a, b *Item[T]) bool) Option[T] {
	return func(q *PriorityQueue[T]) {
		q.lessFunc = less
	}
}

func WithEqualFunc[T any](equal func(a, b T) bool) Option[T] {
	return func(q *PriorityQueue[T]) {
		q.equalFunc = equal
	}
}

func NewPriorityQueue[T any](opts ...Option[T]) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		priorityQueue: priorityQueue[T]{
			lessFunc:  func(a, b *Item[T]) bool { return a.priority > b.priority },
			equalFunc: func(a, b T) bool { return reflect.DeepEqual(a, b) },
		},
		itemPool: pool.New(func() *Item[T] {
			return new(Item[T])
		}),
	}
	for _, opt := range opts {
		opt(pq)
	}
	return pq
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

// Iterator returns all elements in the queue sorted by priority (highest first).
func (pq *PriorityQueue[T]) Iterator() []T {

	pq.mu.RLock()
	defer pq.mu.RUnlock()

	itemsCopy := make([]*Item[T], len(pq.items))
	copy(itemsCopy, pq.items)

	tempQueue := &priorityQueue[T]{items: itemsCopy, lessFunc: pq.lessFunc}
	heap.Init(tempQueue)

	result := make([]T, len(itemsCopy))
	for i := 0; i < len(itemsCopy); i++ {
		result[i] = heap.Pop(tempQueue).(*Item[T]).value
	}

	return result
}

// Remove removes an element matching the provided value and maintains the queue order.
// It returns true if the element was found and removed, otherwise false.
func (pq *PriorityQueue[T]) Remove(data T) bool {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	index := -1
	for i, item := range pq.items {
		if pq.equalFunc(data, item.value) {
			index = i
			break
		}
	}

	if index == -1 {
		return false
	}

	removedItem := heap.Remove(pq, index).(*Item[T])

	removedItem.priority = 0
	removedItem.value = *new(T)
	pq.itemPool.Put(removedItem)
	return true
}
