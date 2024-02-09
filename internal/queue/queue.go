package queue

import "fmt"

type Queue struct {
	arr    []int
	length int
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Arr() []int {
	return q.arr
}

func (q *Queue) Length() int {
	return q.length
}

func (q *Queue) Enqueue(value int) {
	q.arr = append(q.arr, value)
	q.length += 1
}

func (q *Queue) Dequeue() (int, error) {
	if q.length == 0 {
		return 0, fmt.Errorf("Queue is empty")
	}
	value := (q.arr)[0]
	q.arr = (q.arr)[1:]
	q.length -= 1
	return value, nil
}
