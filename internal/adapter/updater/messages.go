package updater

import (
	"LectioBot/internal/adapter/keyboards"
	"LectioBot/internal/adapter/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *Updater) handleMessage() {
	text := u.Update.Message.Text
	u.ChatID = u.Update.Message.Chat.ID

	switch text {
	case "📷 Отправить фото с лекции":
		u.SetStates("waiting_photo")
		msg := tgbotapi.NewMessage(u.ChatID, "Отправьте фото и подпись")
		msg.ReplyMarkup = keyboards.GetCancelKeyBoard()
		u.Ctx.Bot.Send(msg)

	case "⚙️ Конфигурация учебного семестра":
		u.SetStates("check_course_config")
		msg := tgbotapi.NewMessage(u.ChatID, "Выберите параметры")
		msg.ReplyMarkup = keyboards.GetConfigKeyboard()
		u.Ctx.Bot.Send(msg)

	case "👤 Профиль":
		usr := user.NewUserData(u.Ctx, u.ChatID, u.Update)
		usr.GetProfile()

	default:
		u.handleOtherMessages()
	}
}

func (u *Updater) handleOtherMessages() {
	u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Я не понимаю это сообщение"))
}
