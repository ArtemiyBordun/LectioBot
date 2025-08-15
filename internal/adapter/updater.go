package adapter

import (
	"LectioBot/internal/adapter/keyboards"
	"LectioBot/internal/context"
	"LectioBot/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Updater struct {
	Ctx    *context.AppContext
	ChatID int64
	Update tgbotapi.Update
}

func NewUpdater(ctx *context.AppContext) *Updater {
	return &Updater{
		Ctx: ctx,
	}
}

func NewUpdaterForChat(ctx *context.AppContext, chatID int64) *Updater {
	return &Updater{
		Ctx:    ctx,
		ChatID: chatID,
	}
}

func (u *Updater) CheckUpdates() {
	t := tgbotapi.NewUpdate(0)
	t.Timeout = 60
	updates := u.Ctx.Bot.GetUpdatesChan(t)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		u.ChatID = update.Message.Chat.ID
		u.Update = update

		if state, ok := u.Ctx.States[u.ChatID]; ok {
			u.handleState(state)
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				u.startCommand()
			default:
				u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Неизвестная команда"))
			}
			continue
		}

		u.handleMessage()
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
		repository := storage.NewStudentRepo(u.Ctx.DB)
		registered, _ := repository.IsRegistration(u.ChatID)

		if !registered {
			u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Введите свои ФИО (Иванов Иван Иванович) для регистрации"))
			u.SetStates("registration")
		} else {
			u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Здравствуйте!"))
		}
	}
}

func (u *Updater) handleMessage() {
	text := u.Update.Message.Text
	u.ChatID = u.Update.Message.Chat.ID

	switch text {
	case "📷 Отправить фото с лекции":
		u.SetStates("waiting_photo")

		cancelKeyboard := keyboards.GetCancelKeyBoard()

		msg := tgbotapi.NewMessage(u.ChatID, "Отправьте фото и в подписи к нему номер лекции")
		msg.ReplyMarkup = cancelKeyboard
		u.Ctx.Bot.Send(msg)

	case "⚙️ Конфигурация учебного семестра":
		u.SetStates("check_course_config")

		cancelKeyboard := keyboards.GetCancelKeyBoard()

		msg := tgbotapi.NewMessage(u.ChatID, "Отправьте фото и в подписи к нему номер лекции")
		msg.ReplyMarkup = cancelKeyboard
		u.Ctx.Bot.Send(msg)

	//case "👤 Профиль":

	default:
		u.handleOtherMessages()
	}
}

func (u *Updater) handleState(state *context.UserState) {
	switch state.State {
	case "waiting_photo":
		u.GetPhoto()
	case "waiting_confirm_date":
		u.SendDate()

	case "registration":
		student := u.HandleRegistrationName()
		state.Student = student
	case "registration_group":
		u.HandleRegistrationGroup(state.Student)

	case "waiting_mark":
		u.Marking()

	default:
		u.handleOtherMessages()
	}
}

func (u *Updater) handleOtherMessages() {
	u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Я не понимаю это сообщение"))
}

func (u *Updater) SetStates(state string) {
	if u.Ctx.States[u.ChatID] == nil {
		u.Ctx.States[u.ChatID] = &context.UserState{}
	}
	u.Ctx.States[u.ChatID].State = state
}
