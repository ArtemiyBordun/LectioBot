package user

import (
	"log"
	"strings"

	"LectioBot/internal/adapter/keyboards"
	"LectioBot/internal/context"
	"LectioBot/internal/models"
	"LectioBot/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *UserData) HandleRegistrationName() *models.Student {
	fullName := strings.TrimSpace(u.Update.Message.Text)
	group, fullNameSheet := u.Ctx.Sheet.FindStudentGroup(fullName)

	if group == "" {
		sendError("ФИО не найдено, попробуйте ещё раз", u.ChatID, u.Ctx)
		return nil
	}

	student := models.CreateStudent(u.ChatID, fullNameSheet, u.Update.Message.From.UserName, group)

	keyboard := keyboards.GetYesKeyboard()
	msg := tgbotapi.NewMessage(u.ChatID, "Ваша группа: "+group+"?\nЕсли нет, то напишите свою группу в формате СГН3-11Б")
	msg.ReplyMarkup = keyboard
	u.Ctx.Bot.Send(msg)

	u.Ctx.States[u.ChatID].State = "registration_group"
	return student
}

func (u *UserData) HandleRegistrationGroup(student *models.Student) {
	_, err := u.Ctx.Sheet.IsGroup(u.Update.Message.Text)

	if u.Update.Message.Text != "✅ Подтвердить" && err != nil {
		student.Group = u.Update.Message.Text
	}

	repo := storage.NewStudentRepo(u.Ctx.DB)
	if err := repo.Create(student); err != nil {
		log.Println(err)
		sendError("Ошибка, попробуйте ещё раз", u.ChatID, u.Ctx)
		return
	}
	msg := tgbotapi.NewMessage(u.ChatID, "Регистрация завершена!\nВы: "+student.Name+"\nОбучаетесь в группе: "+student.Group)
	keyboard := keyboards.GetStartKeyboard()
	msg.ReplyMarkup = keyboard
	u.Ctx.Bot.Send(msg)
	delete(u.Ctx.States, u.ChatID)
}

func sendError(text string, chatID int64, ctx *context.AppContext) {
	keyboard := keyboards.GetCancelKeyBoard()
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	ctx.Bot.Send(msg)
}
