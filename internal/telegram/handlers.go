package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/maps"
	"schedulerTelegramBot/internal/queue"
	"schedulerTelegramBot/internal/redisRepo"
	"time"
)

const (
	startCommand     = "start"
	registerCommand  = "reg"
	adminCommand     = "admin"
	defaultMessage   = "Unknown command name"
	startMessage     = "Hello! This is scheduler bot, which put users in the virtual schedules! To register in the virtual schedule please click /reg"
	registerMessage  = "Unimplemented"
	basicTextMessage = "It's not allowed here!"
	leaveMessage     = "Unimplemented!"
	adminMessage     = "Unimplemented!"
)

const (
	startState   = "started telegram bot"
	regState     = "registration"
	unknownState = "unknown"
)

type QueueService interface {
	PutInQueue(message *tgbotapi.Message)

	PollFromQueue()

	ReadDataFromQueue()

	GetBackChannel() *chan queue.User

	GetNotificationChannel() *chan queue.User
}

func (b *Bot) handleTextRequests(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, basicTextMessage)
	redis := *b.redisRepo
	state := redis.GetSession(message.Chat.ID).State
	switch state {
	case startState:
		msg.Text = startMessage
	case regState:
		msg.Text = registerMessage
		var service = *b.s
		go service.PutInQueue(message)
		var channel = service.GetBackChannel()
		go b.ReadFromQueue(channel)
	default:
		msg.Text = basicTextMessage
	}
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleCommandRequests(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, defaultMessage)
	var keyboardBuilder = *b.keyboardBuilder
	delegatedMessage := make(map[int64]string)
	switch message.Command() {
	case startCommand:
		msg.Text = startMessage
		delegatedMessage[message.Chat.ID] = "started telegram bot"
		*b.channel <- delegatedMessage
	case registerCommand:
		msg.Text = registerMessage
		delegatedMessage[message.Chat.ID] = "registration"
		*b.channel <- delegatedMessage
		keyboardBuilder.BuildKeyboard(&msg, []string{"oijuhg", "gyhuji", "YGHUj"})
	default:
		delegatedMessage[message.Chat.ID] = unknownState
		*b.channel <- delegatedMessage
	}
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) Read(ch chan map[int64]string) {
	for {
		select {
		case v := <-ch:
			redis := *b.redisRepo
			var keys = maps.Keys(v)
			for _, val := range keys {
				var session = &redisRepo.Session{ChatId: val, State: v[val]}
				go redis.StartReading(val, session)
			}
		default:
			time.Sleep(3 * time.Second)
		}
	}
}

func (b *Bot) ReadFromQueue(channel *chan queue.User) {
	for {
		select {
		case v := <-*channel:
			msg := tgbotapi.NewMessage(v.ChatId, v.Message)
			b.bot.Send(msg)
		default:
			time.Sleep(3 * time.Second)
		}
	}
}
