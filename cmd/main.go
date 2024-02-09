package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"schedulerTelegramBot/configs"
	"schedulerTelegramBot/internal/redisRepo"
	b "schedulerTelegramBot/internal/telegram"
)

func main() {
	c := context.Background()
	configs := configs.GetInstance()
	redisConfigs := configs.Redis
	var channel = make(chan string)
	redisDB := redisRepo.NewRedisDB(&c, &channel, redisConfigs)
	go redisDB.StartReading()
	bot, err := tgbotapi.NewBotAPI(configs.Bot.Token)
	if err != nil {
		log.Fatal(err)
	}
	v := make(chan string)
	go b.Read(v)
	bot.Debug = true
	telegramBot := b.NewBot(bot, &v)
	err = telegramBot.StartBot()
	if err != nil {
		log.Fatal(err)
	}
}
