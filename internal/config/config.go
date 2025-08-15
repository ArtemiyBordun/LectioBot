package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"LectioBot/internal/models"

	"github.com/joho/godotenv"
)

func LoadConfig() *models.Config {
	err := godotenv.Load("internal/config/config.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки config.env: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN не задан в config.env")
	}

	adminIDsStr := os.Getenv("ADMIN_ID")
	if adminIDsStr == "" {
		log.Fatal("ADMIN_ID не задан в config.env")
	}

	var adminIDs []int64
	for _, idStr := range strings.Split(adminIDsStr, ",") {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Fatalf("Некорректный ADMIN_ID: %v", err)
		}
		adminIDs = append(adminIDs, id)
	}

	lecturerIDStr := os.Getenv("LECTURER_ID")
	lecturerID, err := strconv.ParseInt(lecturerIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Некорректный LECTURER_ID: %v", err)
	}

	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	if spreadsheetID == "" {
		log.Fatal("SPREADSHEET_ID не задан")
	}

	credentialsFile := os.Getenv("CREDENTIALS_FILE")
	if credentialsFile == "" {
		log.Fatal("CREDENTIALS_FILE не задан")
	}

	course, _ := LoadCourseConfig("internal/config/course_config.json")

	return models.GetConfig(botToken, adminIDs, course, lecturerID, spreadsheetID, credentialsFile)
}
