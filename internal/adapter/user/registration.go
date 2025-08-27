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
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "ФИО не найдено, попробуйте ещё раз"))
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
	groupName := u.Update.Message.Text

	if groupName != "✅ Подтвердить" {
		exists, err := u.Ctx.Sheet.IsGroup(groupName)
		if err != nil {
			log.Println("Ошибка проверки группы:", err)
			u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Ошибка при проверке группы, попробуйте ещё раз"))
			return
		}

		if !exists {
			u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Группа не найдена. Введите правильное название группы:"))
			return
		}

		student.Group = groupName

		keyboard := keyboards.GetYesKeyboard()
		msg := tgbotapi.NewMessage(u.ChatID, "Группа '"+groupName+"' найдена. Подтвердите выбор или введите другую группу:")
		msg.ReplyMarkup = keyboard
		u.Ctx.Bot.Send(msg)
		return
	}

	repo := storage.NewStudentRepo(u.Ctx.DB)
	if err := repo.Create(student); err != nil {
		log.Println("Ошибка создания студента:", err)
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Ошибка регистрации, попробуйте ещё раз"))
		return
	}

	msg := tgbotapi.NewMessage(u.ChatID, "Регистрация завершена!\nВы: "+student.Name+"\nОбучаетесь в группе: "+student.Group)
	msg.ReplyMarkup = keyboards.GetStartKeyboard()
	u.Ctx.Bot.Send(msg)
	delete(u.Ctx.States, u.ChatID)
}

func sendError(text string, chatID int64, ctx *context.AppContext) {
	keyboard := keyboards.GetCancelKeyBoard()
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	ctx.Bot.Send(msg)
}
