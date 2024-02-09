package services

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"schedulerTelegramBot/internal/queue"
	"time"
)

type QueueService struct {
	Queue       *queue.Queue
	Channel     *chan queue.User
	BackChannel *chan string
}

func (s *QueueService) PollFromQueue() {
	for {
		select {
		case v := <-*s.Channel:
			fmt.Println(v)
			time.Sleep(10 * time.Second)
			*s.BackChannel <- "backLog"
		default:
			time.Sleep(3 * time.Second)
		}
	}
}

func (s *QueueService) GetBackChannel() *chan string {
	return s.BackChannel
}

func (s *QueueService) PutInQueue(message *tgbotapi.Message) {
	var user = queue.User{Username: message.Chat.UserName, ChatId: message.Chat.ID, PlaceInQueue: s.Queue.Length()}
	*s.Channel <- user
}
