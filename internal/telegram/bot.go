package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"schedulerTelegramBot/internal/redisRepo"
)

type RedisRepo interface {
	StartReading(key int64, session *redisRepo.Session)
}

type Bot struct {
	bot       *tgbotapi.BotAPI
	s         *QueueService
	redisRepo *RedisRepo
	channel   *chan map[int64]string
}

func NewBot(bot *tgbotapi.BotAPI, ch *chan map[int64]string, redis *RedisRepo) *Bot {
	return &Bot{bot: bot, redisRepo: redis, channel: ch}
}

func (b *Bot) StartBot() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}
	b.handleUpdates(updates)
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil { // If we got a message
			if update.Message.IsCommand() {
				err := b.handleCommandRequests(update.Message)
				if err != nil {
					return
				}
			} else {
				err := b.handleTextRequests(update.Message)
				if err != nil {
					return
				}
			}
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			//
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//
			//b.bot.Send(msg)
		} else {
			continue
		}
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u), nil
}
