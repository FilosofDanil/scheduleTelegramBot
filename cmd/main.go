package main

import (
	"schedulerTelegramBot/configs"
	"schedulerTelegramBot/pkg/telegram/bot"
)

func main() {
	configs := configs.GetInstance()
	bot.RegisterBot(&configs.Bot)
}
