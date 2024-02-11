package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"schedulerTelegramBot/internal/queue"
	"schedulerTelegramBot/internal/redisRepo"
	"time"
)

const (
	startCommand                        = "start"
	registerCommand                     = "reg"
	deleteFromQueue                     = "logout"
	adminCommand                        = "admin"
	defaultMessage                      = "Unknown command name"
	startMessage                        = "Hello! This is scheduler bot, which put users in the virtual schedules! To register in the virtual schedule please click /reg"
	registerMessage                     = "Hello! Are you sure about putting yourself in the queue, please confirm it in the provided tab below"
	basicTextMessage                    = "It's not allowed here!"
	leaveMessage                        = "Okay! You still can logout from queue in any time you prefer!"
	adminMessage                        = "Are you sure, that you would like to log in into admin tab?"
	registrationNoOption                = "Okay! You still can register in the queue in every convenient time!"
	registrationYesOption               = "Okay! Wait for server answer..."
	deleteFromQueueButNotPresentMessage = "You haven't login in any queue yet!"
	deleteFromQueueMessage              = "Are you sure, that you would like to leave the queue(this action cannot be canceled), please confirm it in the provided tab below"
	deleteFromQueueSuccessMessage       = "Successfully deleted from queue. You always have an opportunity to login in the queue in any time!"
	adminYesOptionMessage               = "Okay, the verification id has sent on your commerce email, please check oit and send back in the chat!"
	adminNoOptionMessage                = "Okay, You still can log in admin tab in every convenient time!"
	adminSuccessfully                   = "Congratulations! You have successfully authorized!"
	adminUnSuccessfully                 = "Given uuid is wrong!"
)

const (
	startState            = "started telegram bot"
	regState              = "registration"
	unknownState          = "unknown"
	logoutState           = "logout"
	adminState            = "admin"
	adminWaitForUIIDState = "adminWait"
)

var (
	optionsYesNoKeyboard = []string{"Yes", "No"}
)

type QueueService interface {
	PutInQueue(message *tgbotapi.Message)

	PollFromQueue()

	ReadDataFromQueue()

	DeleteFromQueue(message *tgbotapi.Message) error

	GetBackChannel() *chan queue.User

	GetNotificationChannel() *chan queue.User
}

type EmailService interface {
	SendMail() error

	CheckUUID(idFromMessage string) error
}

func (b *Bot) handleTextRequests(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, basicTextMessage)
	redis := *b.redisRepo
	state := redis.GetSession(message.Chat.ID).State
	var emailService = *b.emailService
	switch state {
	case startState:
		msg.Text = startMessage
	case adminState:
		switch message.Text {
		case "Yes":
			msg.Text = adminYesOptionMessage
			go func() {
				err := emailService.SendMail()
				if err != nil {
					zap.L().Error(err.Error())
				}
			}()
			delegatedMessage := make(map[int64]string)
			delegatedMessage[message.Chat.ID] = adminWaitForUIIDState
			*b.channel <- delegatedMessage
		case "No":
			msg.Text = adminNoOptionMessage
		}
	case adminWaitForUIIDState:
		err := emailService.CheckUUID(message.Text)
		if err != nil {
			zap.L().Error(err.Error())
			msg.Text = adminUnSuccessfully
		} else {
			msg.Text = adminSuccessfully
		}
	case regState:
		switch message.Text {
		case "Yes":
			var service = *b.queueService
			go service.PutInQueue(message)
			var channel = service.GetBackChannel()
			go b.ReadFromQueue(channel)
			go b.ReadFromQueue(service.GetNotificationChannel())
			msg.Text = registrationYesOption
		case "No":
			msg.Text = registrationNoOption
		}
	case logoutState:
		switch message.Text {
		case "Yes":
			var service = *b.queueService
			err := service.DeleteFromQueue(message)
			if err != nil {
				msg.Text = deleteFromQueueButNotPresentMessage
			} else {
				msg.Text = deleteFromQueueSuccessMessage
			}
		case "No":
			msg.Text = leaveMessage
		}
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
		delegatedMessage[message.Chat.ID] = startState
	case registerCommand:
		msg.Text = registerMessage
		delegatedMessage[message.Chat.ID] = regState
		keyboardBuilder.BuildKeyboard(&msg, optionsYesNoKeyboard)
	case deleteFromQueue:
		msg.Text = deleteFromQueueMessage
		delegatedMessage[message.Chat.ID] = logoutState
	case adminCommand:
		msg.Text = adminMessage
		keyboardBuilder.BuildKeyboard(&msg, optionsYesNoKeyboard)
		delegatedMessage[message.Chat.ID] = adminState
	default:
		delegatedMessage[message.Chat.ID] = unknownState
	}
	*b.channel <- delegatedMessage
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
			_, err := b.bot.Send(msg)
			if err != nil {
				zap.L().Error(err.Error())
			}
		default:
			time.Sleep(3 * time.Second)
		}
	}
}
