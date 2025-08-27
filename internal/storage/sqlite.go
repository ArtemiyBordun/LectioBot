package storage

import (
	"LectioBot/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitSQLite(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Student{}, &models.Lecture{}, &models.Attendance{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
