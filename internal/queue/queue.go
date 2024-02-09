package queue

import (
	"fmt"
)

type User struct {
	ChatId       string
	Username     string
	PlaceInQueue int
}

type Queue struct {
	arr    []User
	length int
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Arr() []User {
	return q.arr
}

func (q *Queue) Length() int {
	return q.length
}

func (q *Queue) Enqueue(value User) {
	q.arr = append(q.arr, value)
	q.length += 1
}

func (q *Queue) Dequeue() (User, error) {
	if q.length == 0 {
		return User{}, fmt.Errorf("Queue is empty")
	}
	value := (q.arr)[0]
	q.arr = (q.arr)[1:]
	q.length -= 1
	return value, nil
}
