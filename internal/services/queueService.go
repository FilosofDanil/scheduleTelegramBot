package services

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"schedulerTelegramBot/internal/queue"
	"time"
)

type QueueService struct {
	Queue               *queue.Queue
	Channel             *chan queue.User
	BackChannel         *chan queue.User
	NotificationChannel *chan queue.User
}

func (s *QueueService) PollFromQueue() {
	for {
		select {
		case v := <-*s.Channel:
			fmt.Println(v)
			time.Sleep(10 * time.Second)
			v.Message = "backLog: " + v.Username
			*s.BackChannel <- v
		default:
			time.Sleep(3 * time.Second)
		}
	}
}

func (s *QueueService) GetBackChannel() *chan queue.User {
	return s.BackChannel
}

func (s *QueueService) GetNotificationChannel() *chan queue.User {
	return s.NotificationChannel
}

func (s *QueueService) PutInQueue(message *tgbotapi.Message) {
	var user = queue.User{Username: message.Chat.UserName, ChatId: message.Chat.ID, PlaceInQueue: s.Queue.Length()}
	s.Queue.Enqueue(user)
	*s.Channel <- user
}

func (s *QueueService) ReadDataFromQueue() {
	for {
		time.Sleep(30 * time.Second)
		var user, _ = s.Queue.Dequeue()
		fmt.Println("test")
		fmt.Println(user)
		user.Message = "Your time!"
		*s.NotificationChannel <- user
	}
}
