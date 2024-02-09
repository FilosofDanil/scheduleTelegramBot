package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type KeyBoardBuilderService struct {
	settings map[string]string
}

func NewKeyBoardBuilderService(settings map[string]string) *KeyBoardBuilderService {
	return &KeyBoardBuilderService{settings: settings}
}

func (k KeyBoardBuilderService) BuildKeyboard(message *tgbotapi.MessageConfig, rows [][]string) {
	var button = tgbotapi.NewKeyboardButton("wewrggerre")
	message.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button})
}
