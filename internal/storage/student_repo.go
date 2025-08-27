package storage

import (
	"LectioBot/internal/models"

	"gorm.io/gorm"
)

type StudentRepo struct {
	db *gorm.DB
}

func NewStudentRepo(db *gorm.DB) *StudentRepo {
	return &StudentRepo{db: db}
}

func (r *StudentRepo) Create(student *models.Student) error {
	return r.db.Create(student).Error
}

func (r *StudentRepo) GetStudentByChatID(chatID int64) (*models.Student, error) {
	var student models.Student
	if err := r.db.Where("chat_id = ?", chatID).First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *StudentRepo) GetAll() ([]models.Student, error) {
	var students []models.Student
	if err := r.db.Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

func (r *StudentRepo) GetAllIDs() ([]int64, error) {
	var ids []int64
	if err := r.db.Model(&models.Student{}).Pluck("chat_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *StudentRepo) IsRegistration(chatId int64) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Student{}).
		Where("chat_id = ?", chatId).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
