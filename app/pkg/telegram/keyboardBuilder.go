package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type KeyBoardBuilderService struct {
	settings map[string]string
}

func NewKeyBoardBuilderService(settings map[string]string) *KeyBoardBuilderService {
	return &KeyBoardBuilderService{settings: settings}
}

func (k *KeyBoardBuilderService) BuildKeyboard(message *tgbotapi.MessageConfig, rows []string) {
	var buttons []tgbotapi.KeyboardButton

	for _, row := range rows {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(row))
	}
	var keyboard = tgbotapi.NewReplyKeyboard(buttons)
	keyboard.OneTimeKeyboard = true
	message.ReplyMarkup = keyboard

	//TODO add implementation for InlineKeyboardMarkup
}

func (k *KeyBoardBuilderService) ManageSettings(key string, value string) {
	//TODO implement keyBoard management
}
