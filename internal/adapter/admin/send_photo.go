package admin

import (
	"LectioBot/internal/adapter/keyboards"
	"LectioBot/internal/adapter/user"
	"LectioBot/internal/context"
	"LectioBot/internal/models"
	"LectioBot/internal/storage"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewUserForChat(ctx *context.AppContext, chatID int64) *user.UserData {
	return &user.UserData{
		Ctx:    ctx,
		ChatID: chatID,
	}
}

func (u *AdminData) GetPhoto() {
	if u.Update.Message.Text == "❌ Отмена" {
		delete(u.Ctx.States, u.ChatID)

		var keyboard tgbotapi.ReplyKeyboardMarkup
		if u.Ctx.Cfg.IsAdmin(u.ChatID) {
			keyboard = keyboards.GetAdminKeyBoard()
		} else {
			keyboard = keyboards.GetStartKeyboard()
		}
		msg := tgbotapi.NewMessage(u.ChatID, "Действие отменено")
		msg.ReplyMarkup = keyboard
		u.Ctx.Bot.Send(msg)
		return
	}
	lectureNumber := u.Update.Message.Caption
	if lectureNumber == "" {
		keyboard := keyboards.GetAdminKeyBoard()
		msg := tgbotapi.NewMessage(u.ChatID, "Вы не ввели номер лекции 😭")
		msg.ReplyMarkup = keyboard
		u.Ctx.Bot.Send(msg)
		return
	}

	num, err := strconv.Atoi(lectureNumber)
	if err != nil {
		keyboard := keyboards.GetAdminKeyBoard()
		msg := tgbotapi.NewMessage(u.ChatID, "а можно нормальные циферки")
		msg.ReplyMarkup = keyboard
		u.Ctx.Bot.Send(msg)
		return
	}

	if num <= u.Ctx.LastLecture {
		keyboard := keyboards.GetAdminKeyBoard()
		msg := tgbotapi.NewMessage(u.ChatID, "а можно нормальные циферки, больше чем предыдущие")
		msg.ReplyMarkup = keyboard
		u.Ctx.Bot.Send(msg)
		return
	}

	if len(u.Update.Message.Photo) > 0 {
		photos := u.Update.Message.Photo
		fileID := photos[len(photos)-1].FileID
		u.sendDone(fileID, "image", lectureNumber)

	} else if u.Update.Message.Document != nil && isImage(u.Update.Message.Document.MimeType) {
		fileID := u.Update.Message.Document.FileID
		u.sendDone(fileID, "doc", lectureNumber)
	} else {
		u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Это не фото или файл с фото. Попробуйте ещё раз."))
	}
}

func (u *AdminData) sendPhoto(fileID, fileType, lectureNumber, lectureDate string) {
	repo := storage.NewStudentRepo(u.Ctx.DB)
	ids, err := repo.GetAllIDs()
	if err == nil {
		for _, id := range ids {
			u.sendPhotoToUser(id, fileID, fileType, lectureNumber, lectureDate)
		}
	}
}

func (u *AdminData) sendPhotoToUser(chatID int64, fileID, fileType, lectureNumber, lectureDate string) {
	updater := NewUserForChat(u.Ctx, chatID)

	caption := "Время отметиться на лекции номер " + lectureNumber + ", которая была " + lectureDate +
		"\nНайди себя на фото и отправь сюда свой номер 👇"

	var msg tgbotapi.Chattable
	var keyboard tgbotapi.ReplyKeyboardMarkup
	if !u.Ctx.Cfg.IsAdmin(chatID) {
		keyboard = keyboards.GetLectureKeyboard()
	} else {
		keyboard = keyboards.GetSendKeyboard()
	}

	switch fileType {
	case "image":
		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(fileID))
		photo.Caption = caption
		photo.ReplyMarkup = keyboard
		msg = photo

	case "doc":
		doc := tgbotapi.NewDocument(chatID, tgbotapi.FileID(fileID))
		doc.Caption = caption
		doc.ReplyMarkup = keyboard
		msg = doc
	}
	updater.Ctx.Bot.Send(msg)

	if !u.Ctx.Cfg.IsAdmin(chatID) {
		updater.SetStates("waiting_mark")
	}
}

func isImage(mime string) bool {
	return mime == "image/jpeg" || mime == "image/png" || mime == "image/jpg"
}

func (u *AdminData) sendDone(fileID, fileType, lectureNumberStr string) {
	keyboard := keyboards.GetConfirmKeyboard()
	msg := tgbotapi.NewMessage(u.ChatID, "Файл с фото получен!\nЛекция была сегодня? Если нет, то напишите её дату в формате дд.мм.гггг")
	msg.ReplyMarkup = keyboard
	u.Ctx.Bot.Send(msg)
	u.Ctx.States[u.ChatID].State = "waiting_confirm_date"
	u.Ctx.PhotoData = context.NewPhotoData(fileID, fileType, lectureNumberStr)
}

func (u *AdminData) SendDate() string {
	var data time.Time
	var dateStr string

	if u.Update.Message.Text == "Да" {
		now := time.Now()
		data = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		dateStr = data.Format("02.01.2006")
	} else {
		dateStr = u.Update.Message.Text
		parsed, err := time.Parse("02.01.2006", dateStr)
		if err != nil {
			u.Ctx.Bot.Send(tgbotapi.NewMessage(u.ChatID, "Неверный формат даты, используйте дд.мм.гггг"))
			return ""
		}
		data = parsed
	}
	u.Date = dateStr
	u.Ctx.States[u.ChatID].State = "waiting_confirm"

	u.sendPhotoToUser(u.ChatID, u.Ctx.PhotoData.FileID, u.Ctx.PhotoData.FileType, u.Ctx.PhotoData.LectureNumberStr, u.Date)

	return dateStr
}

func (u *AdminData) Send() {
	if u.Update.Message.Text == "Отправить" {
		keyboard := keyboards.GetAdminKeyBoard()
		msg := tgbotapi.NewMessage(u.ChatID, "Отлично, фото отправлено")
		msg.ReplyMarkup = keyboard
		u.Ctx.Bot.Send(msg)

		lectureNumber, _ := strconv.Atoi(u.Ctx.PhotoData.LectureNumberStr)
		u.Ctx.LastLecture = lectureNumber
		u.sendPhoto(u.Ctx.PhotoData.FileID, u.Ctx.PhotoData.FileType, u.Ctx.PhotoData.LectureNumberStr, u.Date)
		delete(u.Ctx.States, u.ChatID)

		repo := storage.NewLectureRepo(u.Ctx.DB)
		repo.Create(models.CreateLecture(lectureNumber, 0, u.Date))
	} else {
		keyboard := keyboards.GetAdminKeyBoard()
		msg := tgbotapi.NewMessage(u.ChatID, "отправка отменена")
		msg.ReplyMarkup = keyboard
		u.Ctx.Bot.Send(msg)
	}
}
