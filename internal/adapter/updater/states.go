package updater

import (
	"LectioBot/internal/adapter/keyboards"
	"LectioBot/internal/adapter/user"
	"LectioBot/internal/context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *Updater) handleState(sub *SubUpdaters, state *context.UserState) {
	switch state.State {
	case "waiting_photo":
		sub.AdminData.GetPhoto()
	case "waiting_confirm_date":
		sub.AdminData.Date = sub.AdminData.SendDate()
	case "waiting_confirm":
		sub.AdminData.Send()

	case "registration":
		usr := user.NewUserData(u.Ctx, u.ChatID, u.Update)
		student := usr.HandleRegistrationName()
		state.Student = student

	case "registration_group":
		usr := user.NewUserData(u.Ctx, u.ChatID, u.Update)
		usr.HandleRegistrationGroup(state.Student)

	case "waiting_mark":
		usr := user.NewUserData(u.Ctx, u.ChatID, u.Update)
		usr.Marking()

	case "check_course_config":
		if u.Update.Message.Text == "❌ Отмена" {
			// отмена
		}
		msg := tgbotapi.NewMessage(u.ChatID, "help")
		msg.ReplyMarkup = keyboards.GetAdminKeyBoard()
		u.Ctx.Bot.Send(msg)
		delete(u.Ctx.States, u.ChatID)

	default:
		u.handleOtherMessages()
	}
}
