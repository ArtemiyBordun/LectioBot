package user

import (
	"LectioBot/internal/adapter/keyboards"
	"LectioBot/internal/models"
	"LectioBot/internal/storage"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *UserData) Marking() {
	if u.Update.Message.Text != "Меня не было(((" {
		num, err := strconv.Atoi(u.Update.Message.Text)
		if err != nil || num < 1 || num > len(u.Ctx.Sheet.Students) {
			u.sendAnswer("Это не номер, попробуй ещё раз")
			return
		}
		repo_attendance := storage.NewAttendanceRepo(u.Ctx.DB)
		isUnique, err := repo_attendance.GetUniqueNumberByLection(u.Ctx.LastLecture)

		if err != nil {
			u.sendAnswer("Произошла ошибка при проверке номера, попробуйте ещё раз")
			return
		}

		if !isUnique {
			u.sendAnswer("Такой номер уже занят, если возникли проблемы обратитесь к семинаристу/лектору")
			delete(u.Ctx.States, u.ChatID)
			return
		}

		repo_student := storage.NewStudentRepo(u.Ctx.DB)
		student, err := repo_student.GetStudentByChatID(u.ChatID)
		if err == nil {
			res := u.Ctx.Sheet.Marking(student, u.Ctx.LastLecture)
			if res == "" {
				u.updateCountStudent(u.Update.Message.Text)
				u.sendAnswer("Отлично! Посещение отмечено")
			} else {
				u.sendAnswer(res + "\nЕсли возникли проблемы, обращайтесь к семинаристу/лектору")
			}
		}
	} else {
		u.sendAnswer("Плохо")
	}
	delete(u.Ctx.States, u.ChatID)
}

func (u *UserData) sendAnswer(text string) {
	keyboard := keyboards.GetStartKeyboard()
	msg := tgbotapi.NewMessage(u.ChatID, text)
	msg.ReplyMarkup = keyboard
	u.Ctx.Bot.Send(msg)
}

func (u *UserData) updateCountStudent(num string) {
	repoLecture := storage.NewLectureRepo(u.Ctx.DB)
	repoLecture.IncrementCountStudent(u.Ctx.LastLecture)

	number, _ := strconv.Atoi(num)

	repoStudent := storage.NewStudentRepo(u.Ctx.DB)
	student, _ := repoStudent.GetStudentByChatID(u.ChatID)

	repoAttendance := storage.NewAttendanceRepo(u.Ctx.DB)
	attendance := models.CreateAttendance(u.Ctx.LastLecture, number, student.ChatId, student.Name, student.Group)
	repoAttendance.Create(attendance)
}
