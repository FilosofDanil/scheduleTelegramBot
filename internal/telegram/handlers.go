package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	startCommand     = "start"
	registerCommand  = "reg"
	leave            = "leave"
	defaultMessage   = "Unknown command name"
	startMessage     = "Hello! This is scheduler bot, which put users in the virtual schedules! To register in the virtual schedule please click /reg"
	registerMessage  = "Unimplemented"
	basicTextMessage = "It's not allowed here!"
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
	}
	_, err := b.bot.Send(msg)
	return err
}
