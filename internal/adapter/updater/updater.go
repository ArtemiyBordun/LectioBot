package updater

import (
	"LectioBot/internal/adapter/admin"
	"LectioBot/internal/context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Updater struct {
	Ctx    *context.AppContext
	ChatID int64
	Update tgbotapi.Update
}

type SubUpdaters struct {
	AdminData admin.AdminData
}

func NewUpdater(ctx *context.AppContext) *Updater {
	return &Updater{
		Ctx: ctx,
	}
}

func (u *Updater) NewSubUpdater() *SubUpdaters {
	return &SubUpdaters{
		AdminData: *admin.NewAdminData(u.Ctx, u.ChatID, u.Update),
	}
}

func (u *Updater) SetStates(state string) {
	if u.Ctx.States[u.ChatID] == nil {
		u.Ctx.States[u.ChatID] = &context.UserState{}
	}
	u.Ctx.States[u.ChatID].State = state
}
