package user

import (
	"LectioBot/internal/adapter/keyboards"
	"LectioBot/internal/models"
	"LectioBot/internal/storage"
	"LectioBot/pkg/external"

	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *UserData) GetProfile() {
	student, err := u.getStudent()
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(u.ChatID, "")
	msg.ParseMode = "HTML"

	msg.Text = u.buildProfileText(student)
	msg.ReplyMarkup = keyboards.GetStartKeyboard()

	u.Ctx.Bot.Send(msg)
}

func (u *UserData) getStudent() (*models.Student, error) {
	studentRepo := storage.NewStudentRepo(u.Ctx.DB)
	student, err := studentRepo.GetStudentByChatID(u.ChatID)
	if err != nil {
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Ошибка при получении профиля"))
		return nil, err
	}
	return student, nil
}

func (u *UserData) buildProfileText(student *models.Student) string {
	text := fmt.Sprintf(
		"Ваш профиль:\n\nВаш Chat_id: <code>%d</code>\nИмя: <i>%s</i>\nГруппа: <i>%s</i>\n",
		student.ChatId, student.Name, student.Group,
	)

	studentSheet := u.Ctx.Sheet.GetStudentByName(student.Name)
	if studentSheet != nil {
		text += u.buildPointsText(studentSheet)
	}

	attendanceCount := u.getAttendance(student.ChatId)
	text += "\nКоличество посещений: <i>" + strconv.Itoa(attendanceCount) + "</i> из " + strconv.Itoa(u.Ctx.Cfg.CourseConfig.TotalLectures) + "\n"

	text += "\n<tg-spoiler>Если есть какие-то вопросы или проблемы обращайтесь лично/через чат к своему <i>семинаристу</i> или <i>лектору</i></tg-spoiler>\n"

	return text
}

func (u *UserData) buildPointsText(sheet *external.Student) string {
	text := fmt.Sprintf(
		"\n📊 Баллы по ОП: %d/%d\n",
		int(sheet.PointsOP), int(u.Ctx.Cfg.CourseConfig.MaxPoints),
	)

	remainingOP := int(u.Ctx.Cfg.CourseConfig.MinPointsOPForPass) - int(sheet.PointsOP)
	if remainingOP <= 0 {
		text += "✅ Допуск к зачету/экзамену получен!\n"
	} else {
		text += fmt.Sprintf("⚠️ Для допуска к зачету/экзамену необходимо ещё <i>%d</i> баллов\n", remainingOP)
	}

	if sheet.PointsOOP != -1 {
		text += fmt.Sprintf(
			"📊 Баллы по ООП: %d/%d\n",
			int(sheet.PointsOOP), int(u.Ctx.Cfg.CourseConfig.MaxPoints),
		)

		remainingOOP := int(u.Ctx.Cfg.CourseConfig.MinPointsOOPForPass) - int(sheet.PointsOOP)
		if remainingOOP <= 0 {
			text += "✅ Допуск к экзамену по ООП получен!\n"
		} else {
			text += fmt.Sprintf("⚠️ Для допуска к экзамену по ООП необходимо ещё <i>%d</i> баллов\n", remainingOOP)
		}
	}

	return text
}

func (u *UserData) getAttendance(chatID int64) int {
	attendanceRepo := storage.NewAttendanceRepo(u.Ctx.DB)
	count, err := attendanceRepo.GetStudentAttendanceCount(chatID)
	if err != nil {
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Ошибка при получении посещаемости"))
		return 0
	}
	return int(count)
}

func (u *UserData) GetHistory() {
	lectureRepo := storage.NewLectureRepo(u.Ctx.DB)
	attendanceRepo := storage.NewAttendanceRepo(u.Ctx.DB)

	lectures, err := lectureRepo.GetAll()
	if err != nil {
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Ошибка при получении лекций"))
		return
	}

	records, err := attendanceRepo.GetStudentAttendance(u.ChatID)
	if err != nil {
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Ошибка при получении посещений"))
		return
	}

	attended := make(map[int]bool)
	for _, rec := range records {
		attended[rec.LectureId] = true
	}

	history := "🕒 История посещений:\n\n"
	for i, l := range lectures {
		status := "❌"
		if attended[l.Id] {
			status = "✅"
		}
		history += fmt.Sprintf("Лекция %d (%s) – %s\n", i+1, l.Date, status)
	}

	u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, history))
}
