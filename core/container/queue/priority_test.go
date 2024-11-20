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

	fmt.Println(pq.Dequeue()) // Task 2 true
	fmt.Println(pq.Dequeue()) // Task 3 true
	fmt.Println(pq.Dequeue()) // Task 1 true

	// Queue is empty
	if value, ok := pq.Dequeue(); ok {
		fmt.Println("Dequeued", value)
	} else {
		fmt.Println("Queue is empty")
	}
}
