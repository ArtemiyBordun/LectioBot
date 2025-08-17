package admin

import (
	"LectioBot/internal/context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdminData struct {
	Ctx    *context.AppContext
	ChatID int64
	Update tgbotapi.Update
}

func NewAdminData(ctx *context.AppContext, chatId int64, update tgbotapi.Update) *AdminData {
	return &AdminData{
		Ctx:    ctx,
		ChatID: chatId,
		Update: update,
	}
}
