package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"schedulerTelegramBot/configs"
	"schedulerTelegramBot/internal/queue"
	"schedulerTelegramBot/internal/redisRepo"
	"schedulerTelegramBot/internal/services"
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
	backChannel := make(chan string)
	userChan := make(chan queue.User)
	bot.Debug = true
	var queueService b.QueueService
	queueService = &services.QueueService{Queue: queue.NewQueue(), Channel: &userChan, BackChannel: &backChannel}
	go queueService.PollFromQueue()
	telegramBot := b.NewBot(bot, &v, &redisDB, &queueService)
	go telegramBot.Read(v)
	err = telegramBot.StartBot()
	if err != nil {
		log.Fatal(err)
	}
}
