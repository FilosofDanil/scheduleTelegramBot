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
	var redisDB b.RedisRepo
	redisDB = redisRepo.NewRedisDB(&c, redisConfigs)
	bot, err := tgbotapi.NewBotAPI(configs.Bot.Token)
	if err != nil {
		log.Fatal(err)
	}
	v := make(chan map[int64]string)

	bot.Debug = true
	telegramBot := b.NewBot(bot, &v, &redisDB)
	go telegramBot.Read(v)
	err = telegramBot.StartBot()
	if err != nil {
		log.Fatal(err)
	}
}
