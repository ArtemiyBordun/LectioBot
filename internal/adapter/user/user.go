package user

import (
	"LectioBot/internal/context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserData struct {
	Ctx    *context.AppContext
	ChatID int64
	Update tgbotapi.Update
}

func NewUserData(ctx *context.AppContext, chatId int64, update tgbotapi.Update) *UserData {
	return &UserData{
		Ctx:    ctx,
		ChatID: chatId,
		Update: update,
	}
}

func (u *UserData) SetStates(state string) {
	if u.Ctx.States[u.ChatID] == nil {
		u.Ctx.States[u.ChatID] = &context.UserState{}
	}
	u.Ctx.States[u.ChatID].State = state
}
