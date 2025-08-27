package updater

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *Updater) CheckUpdates(sub *SubUpdaters) {
	t := tgbotapi.NewUpdate(0)
	t.Timeout = 60
	updates := u.Ctx.Bot.GetUpdatesChan(t)

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		u.ChatID = update.FromChat().ID
		u.Update = update

		if update.Message != nil && update.Message.IsCommand() {
			u.handleCommand()
			continue
		}

		if state, ok := u.Ctx.States[u.ChatID]; ok {
			sub.AdminData.SetUpdateData(u.Ctx, u.ChatID, u.Update)
			u.handleState(sub, state)
			continue
		}

		if update.CallbackQuery != nil {
			u.handleCallback()
			continue
		}

		if update.Message != nil {
			u.handleMessage()
		}
	}
}
