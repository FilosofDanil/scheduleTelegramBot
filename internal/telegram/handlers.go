package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/maps"
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

type QueueService interface {
}

func (b *Bot) handleTextRequests(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, basicTextMessage)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleCommandRequests(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, defaultMessage)
	switch message.Command() {
	case startCommand:
		msg.Text = startMessage
	case registerCommand:
		msg.Text = registerMessage
		delegatedMessage := make(map[int64]string)
		delegatedMessage[message.Chat.ID] = message.Text
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
