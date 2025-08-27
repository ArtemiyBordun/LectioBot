package updater

import (
	"LectioBot/internal/adapter/keyboards"
	"LectioBot/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *Updater) handleCommand() {
	switch u.Update.Message.Command() {
	case "start":
		u.startCommand()
	default:
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Неизвестная команда"))
	}
}

func (u *Updater) startCommand() {
	u.ChatID = u.Update.Message.From.ID

	if u.ChatID == u.Ctx.Cfg.LecturerID || u.Ctx.Cfg.IsAdmin(u.ChatID) {
		adminKeyboard := keyboards.GetAdminKeyBoard()
		var msg tgbotapi.MessageConfig
		if u.ChatID == u.Ctx.Cfg.LecturerID {
			msg = tgbotapi.NewMessage(u.ChatID, "Здравствуйте, Максим Александрович!")
		}
		msg.ReplyMarkup = adminKeyboard
		u.Ctx.Bot.Send(msg)
	} else {
		repo := storage.NewStudentRepo(u.Ctx.DB)
		registered, _ := repo.IsRegistration(u.ChatID)

		if !registered {
			u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Введите свои ФИО для регистрации"))
			u.SetStates("registration")
		} else {
			u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Здравствуйте!"))
		}
	}
}
