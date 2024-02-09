package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"schedulerTelegramBot/internal/redisRepo"
)

type KeyboardBuilder interface {
	ManageSettings(key string, value string)

	BuildKeyboard(message *tgbotapi.MessageConfig, rows []string)
}

type RedisRepo interface {
	StartReading(key int64, session *redisRepo.Session)

	GetSession(key int64) (session *redisRepo.Session)
}

type Bot struct {
	bot             *tgbotapi.BotAPI
	s               *QueueService
	keyboardBuilder *KeyboardBuilder
	redisRepo       *RedisRepo
	channel         *chan map[int64]string
}

func NewBot(bot *tgbotapi.BotAPI, ch *chan map[int64]string, redis *RedisRepo, qs *QueueService) *Bot {
	var keyboardBuilder KeyboardBuilder
	keyboardBuilder = NewKeyBoardBuilderService(make(map[string]string))
	return &Bot{bot: bot, redisRepo: redis, channel: ch, keyboardBuilder: &keyboardBuilder, s: qs}
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
