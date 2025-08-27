package updater

import (
	"LectioBot/internal/adapter/admin"
	"LectioBot/internal/adapter/keyboards"
	"LectioBot/internal/adapter/user"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *Updater) handleMessage() {
	text := u.Update.Message.Text
	u.ChatID = u.Update.Message.Chat.ID

	switch text {
	case "📷 Отправить фото с лекции":
		u.SetStates("waiting_photo")
		str := fmt.Sprintf("Отправьте фото и подписью номер лекции\nПоследняя лекция была под номером %d", u.Ctx.LastLecture)
		msg := tgbotapi.NewMessage(u.ChatID, str)
		msg.ReplyMarkup = keyboards.GetCancelKeyBoard()
		u.Ctx.Bot.Send(msg)

	case "⚙️ Конфигурация учебного семестра":
		msg := tgbotapi.NewMessage(u.ChatID, "Выберите параметры")
		msg.ReplyMarkup = keyboards.GetConfigKeyboard()
		u.Ctx.Bot.Send(msg)

	case "📊 Статистика":
		adm := admin.NewAdminData(u.Ctx, u.ChatID, u.Update)
		adm.GetStatistics()

	case "👤 Профиль":
		usr := user.NewUserData(u.Ctx, u.ChatID, u.Update)
		usr.GetProfile()

	case "🕒 История посещения лекций":
		usr := user.NewUserData(u.Ctx, u.ChatID, u.Update)
		usr.GetHistory()

	default:
		u.handleOtherMessages()
	}
}

func (u *Updater) handleOtherMessages() {
	u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Я не понимаю это сообщение"))
}
