package services

import (
	"fmt"
	"schedulerTelegramBot/internal/queue"
	"time"
)

type QueueService struct {
	Queue   *queue.Queue
	Channel *chan queue.User
}

func (s *QueueService) PutInQueue() {
	for {
		select {
		case v := <-*s.Channel:
			fmt.Println(v)
		default:
			time.Sleep(3 * time.Second)
		}
	}
}

func (s *QueueService) PollFromQueue() {
	var user = queue.User{}
	*s.Channel <- user
}
