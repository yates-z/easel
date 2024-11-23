package queue

import (
	"fmt"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue[string]()
	pq.Enqueue("Task 1", 3)
	pq.Enqueue("Task 2", 1)
	pq.Enqueue("Task 3", 2)

	fmt.Println(pq.Dequeue()) // Task 1 true
	fmt.Println(pq.Dequeue()) // Task 3 true
	fmt.Println(pq.Dequeue()) // Task 2 true

	// Queue is empty
	if value, ok := pq.Dequeue(); ok {
		fmt.Println("Dequeued", value)
	} else {
		fmt.Println("Queue is empty")
	}
}

func Test_Iterator(t *testing.T) {
	pq := NewPriorityQueue[string]()
	pq.Enqueue("Task 1", 1)
	pq.Enqueue("Task 2", 3)
	pq.Enqueue("Task 3", 2)

	//fmt.Printf("%+v %+v %+v\n", pq.items[0], pq.items[1], pq.items[2])

	s := pq.Iterator()
	fmt.Println(s, pq.Len())
	pq.Dequeue()
	fmt.Println(s, pq.Len())
}

func Test_Peek(t *testing.T) {
	pq := NewPriorityQueue[int]()

	pq.Enqueue(10, 1)
	pq.Enqueue(20, 4)
	pq.Enqueue(30, 2)
	pq.Enqueue(40, 3)

	value, ok := pq.Peek()
	fmt.Println("Peek:", value, ok)
}

func Test_Remove(t *testing.T) {
	pq := NewPriorityQueue[int]()

	pq.Enqueue(10, 2)
	pq.Enqueue(20, 1)
	pq.Enqueue(30, 3)

	fmt.Println("Before Remove:", pq.Iterator()) // [30, 10, 20]

	removed := pq.Remove(10)
	fmt.Println("Removed:", removed) // true

	fmt.Println("After Remove:", pq.Iterator()) // [30, 20]

	removed = pq.Remove(40)
	fmt.Println("Removed:", removed) // 输出: false
}
