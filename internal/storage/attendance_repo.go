package storage

import (
	"LectioBot/internal/models"

	"gorm.io/gorm"
)

type AttendanceRepo struct {
	db *gorm.DB
}

func NewAttendanceRepo(db *gorm.DB) *AttendanceRepo {
	return &AttendanceRepo{db: db}
}

func (r *AttendanceRepo) Create(attendance *models.Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *AttendanceRepo) GetMostActiveGroup() (string, error) {
	var result struct {
		Group string
		Count int64
	}

	if err := r.db.Model(&models.Attendance{}).
		Select("`group`, COUNT(*) as count").
		Group("`group`").
		Order("count DESC").
		Limit(1).
		Scan(&result).Error; err != nil {
		return "", err
	}

	if result.Group == "" {
		return "Нет посещений", nil
	}

	return result.Group, nil
}

func (r *AttendanceRepo) GetAverageAttendance() (float64, error) {
	var total int64
	var lecturesCount int64

	if err := r.db.Model(&models.Attendance{}).Count(&total).Error; err != nil {
		return 0, err
	}

	if err := r.db.Model(&models.Attendance{}).Distinct("lecture_id").Count(&lecturesCount).Error; err != nil {
		return 0, err
	}

	if lecturesCount == 0 {
		return 0, nil
	}

	average := float64(total) / float64(lecturesCount)
	return average, nil
}

func (r *AttendanceRepo) GetStudentAttendanceCount(chatID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.Attendance{}).
		Where("chat_id = ?", chatID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *AttendanceRepo) GetStudentAttendance(chatID int64) ([]models.Attendance, error) {
	var records []models.Attendance
	if err := r.db.Where("chat_id = ?", chatID).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
