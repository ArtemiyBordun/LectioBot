package admin

import (
	"LectioBot/internal/context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdminData struct {
	Ctx    *context.AppContext
	ChatID int64
	Update tgbotapi.Update
	Date   string
}

func NewAdminData(ctx *context.AppContext, chatId int64, update tgbotapi.Update) *AdminData {
	return &AdminData{
		Ctx:    ctx,
		ChatID: chatId,
		Update: update,
	}
}

func (u *AdminData) SetUpdateData(ctx *context.AppContext, chatId int64, update tgbotapi.Update) {
	u.Ctx = ctx
	u.ChatID = chatId
	u.Update = update
}
