package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"schedulerTelegramBot/configs"
	b "schedulerTelegramBot/pkg/telegram/bot"
)

func main() {
	configs := configs.GetInstance()
	bot, err := tgbotapi.NewBotAPI(configs.Bot.Token)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	telegramBot := b.NewBot(bot)
	err = telegramBot.StartBot()
	if err != nil {
		log.Fatal(err)
	}
}
