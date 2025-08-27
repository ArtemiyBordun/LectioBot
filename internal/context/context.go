package context

import (
	"LectioBot/internal/models"
	"LectioBot/pkg/external"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type AppContext struct {
	Bot         *tgbotapi.BotAPI
	Cfg         *models.Config
	DB          *gorm.DB
	States      map[int64]*UserState
	Sheet       *external.Sheet
	LastLecture int
	PhotoData   *PhotoData
}

type UserState struct {
	State   string
	Student *models.Student
}

type PhotoData struct {
	FileID           string
	FileType         string
	LectureNumberStr string
}

func NewContext(bot *tgbotapi.BotAPI, cfg *models.Config, db *gorm.DB, sheet *external.Sheet) *AppContext {
	return &AppContext{
		Bot:       bot,
		Cfg:       cfg,
		DB:        db,
		States:    make(map[int64]*UserState),
		Sheet:     sheet,
		PhotoData: NewPhotoData("", "", ""),
	}
}

func NewPhotoData(fileID, fileType, lectureNumberStr string) *PhotoData {
	return &PhotoData{
		FileID:           fileID,
		FileType:         fileType,
		LectureNumberStr: lectureNumberStr,
	}
}
