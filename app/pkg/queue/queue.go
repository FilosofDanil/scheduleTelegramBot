package queue

import (
	"errors"
	"fmt"
)

type User struct {
	ChatId       int64
	Username     string
	PlaceInQueue int
	Message      string
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

func (q *Queue) DeleteFromQueueByChatId(chatId int64) error {
	var i int
	var present = false
	for j, u := range q.arr {
		if u.ChatId == chatId {
			i = j
			present = true
			break
		}
	}
	if !present {
		return errors.New("no value present")
	}
	q.arr[i] = q.arr[len(q.arr)-1]
	q.arr[len(q.arr)-1] = User{}
	q.arr = q.arr[:len(q.arr)-1]
	q.length--
	return nil
}

func (q *Queue) Dequeue() (User, error) {
	if q.length == 0 {
		return User{}, fmt.Errorf("queue is empty")
	}
	value := (q.arr)[0]
	q.arr = (q.arr)[1:]
	q.length -= 1
	return value, nil
}
