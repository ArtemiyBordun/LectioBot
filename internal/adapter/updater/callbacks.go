package updater

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *Updater) handleCallback() {
	data := u.Update.CallbackQuery.Data
	callback := tgbotapi.NewCallback(u.Update.CallbackQuery.ID, "")
	u.Ctx.Bot.Send(callback)

	switch data {
	case "edit_pass_OOP_points":
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Введите новое значение минимальных баллов для допуска к экзамену по ООП"))
		u.SetStates("edit_pass_oop")
	case "edit_pass_OP_points":
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Введите новое значение минимальных баллов для допуска к зачету/экзамену по ОП"))
		u.SetStates("edit_pass_op")
	}
}
