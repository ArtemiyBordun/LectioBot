package admin

import (
	"LectioBot/internal/config"
	"LectioBot/internal/storage"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *AdminData) UpdatePoint(isOP bool) {
	newPointStr := u.Update.Message.Text
	newPoint, err := strconv.Atoi(newPointStr)
	if err != nil {
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "это не число("))
		return
	}
	if newPoint < 0 {
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "ну за что отрицательные баллы..."))
		return
	}
	if !isOP {
		u.Ctx.Cfg.CourseConfig.Update(-1, -1, newPoint, -1)
	} else {
		u.Ctx.Cfg.CourseConfig.Update(-1, newPoint, -1, -1)
	}
	if err := config.SaveCourseConfig(config.CourseConfPath, u.Ctx.Cfg.CourseConfig); err != nil {
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Ошибка сохранения: "+err.Error()))
		return
	}
	u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Баллы успешно сохранены"))
	delete(u.Ctx.States, u.ChatID)
}

func (u *AdminData) GetStatistics() {
	repo := storage.NewLectureRepo(u.Ctx.DB)
	str, _ := repo.GetStatistics(len(u.Ctx.Sheet.Students), storage.NewAttendanceRepo(u.Ctx.DB))
	msg := tgbotapi.NewMessage(u.ChatID, str)
	u.Ctx.Bot.Send(msg)
}
