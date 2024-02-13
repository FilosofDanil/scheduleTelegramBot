package services

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"schedulerTelegramBot/app/pkg/queue"
	"strconv"
	"time"
)

type QueueService struct {
	Queue                        *queue.Queue
	Channel                      *chan queue.User
	BackChannel                  *chan queue.User
	NotificationChannel          *chan queue.User
	SecondaryNotificationChannel *chan queue.User
}

func (s *QueueService) PollFromQueue() {
	for {
		select {
		case v := <-*s.Channel:
			fmt.Println(v)
			time.Sleep(10 * time.Second)
			v.Message = "Congratulations, user: " + v.Username + "\n You're successfully putted in the queue! Now have" +
				" a bit of patience and wait for your turn, your current place in the queue(people after you): " + strconv.Itoa(s.Queue.Length()-1)
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

func (s *QueueService) ReadDataFromQueue() string {
	//for {
	//	time.Sleep(90 * time.Second)
	if s.Queue.Length() == 0 {
		return "The queue is empty"
	}
	var user, _ = s.Queue.Dequeue()
	var usersList = s.Queue.Arr()
	for i, u := range usersList {
		u.Message = "Queue change log, your place in queue(people after you):  " + strconv.Itoa(i)
		*s.SecondaryNotificationChannel <- u
	}
	user.Message = "User: " + user.Username + "\nNow it's your turn in queue! Please follow the further instructions. "
	*s.NotificationChannel <- user
	//}
	return "Chosen user from queue: " + user.Username + "\nCurrent queue length: " + strconv.Itoa(s.Queue.Length())
}

func (s *QueueService) DeleteFromQueue(message *tgbotapi.Message) error {
	err := s.Queue.DeleteFromQueueByChatId(message.Chat.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *QueueService) CheckUser(message *tgbotapi.Message) error {
	for _, u := range s.Queue.Arr() {
		if u.ChatId == message.Chat.ID {
			return errors.New("user already login in the queue")
		}
	}
	return nil
}
