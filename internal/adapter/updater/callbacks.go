package updater

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *Updater) handleCallback() {
	data := u.Update.CallbackQuery.Data
	callback := tgbotapi.NewCallback(u.Update.CallbackQuery.ID, "")
	u.Ctx.Bot.Send(callback)

	switch data {
	case "edit_max_points":
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Введите новое значение максимальных баллов"))
	case "edit_min_points":
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Введите новое значение минимальных баллов"))
	case "cancel":
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Отменено"))
		delete(u.Ctx.States, u.ChatID)
	}
}
