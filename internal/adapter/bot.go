package adapter

import (
	"LectioBot/internal/context"
	"LectioBot/internal/models"
	"LectioBot/pkg/external"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func StartBot(cfg *models.Config, db *gorm.DB, sheet *external.Sheet) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	log.Printf("Бот %s запущен", bot.Self.UserName)

	context := context.NewContext(bot, cfg, db, sheet)

	updater := NewUpdater(context)
	updater.CheckUpdates()
}
