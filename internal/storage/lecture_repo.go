package storage

import (
	"LectioBot/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type LectureRepo struct {
	db *gorm.DB
}

func NewLectureRepo(db *gorm.DB) *LectureRepo {
	return &LectureRepo{db: db}
}

func (r *LectureRepo) Create(lecture *models.Lecture) error {
	return r.db.Create(lecture).Error
}

func (r *LectureRepo) IncrementCountStudent(id int) error {
	return r.db.Model(&models.Lecture{}).
		Where("id = ?", id).
		UpdateColumn("count_student", gorm.Expr("count_student + ?", 1)).Error
}

func (r *LectureRepo) GetAll() ([]models.Lecture, error) {
	var lectures []models.Lecture
	if err := r.db.Find(&lectures).Error; err != nil {
		return nil, err
	}
	return lectures, nil
}

func (r *LectureRepo) GetLectureSummary() (string, error) {
	lectures, err := r.GetAll()
	if err != nil {
		return "", err
	}

	if len(lectures) == 0 {
		return "Лекций ещё не было", nil
	}

	summary := fmt.Sprintf("Всего было лекций %d:\n", len(lectures))
	for i, l := range lectures {
		summary += fmt.Sprintf("%d лекция была %s на ней было %d студентов\n", i+1, l.Date, l.CountStudent)
	}

	return summary, nil
}

func (r *LectureRepo) GetIDByDate(date string) (int, error) {
	var lecture models.Lecture
	if err := r.db.Where("date = ?", date).First(&lecture).Error; err != nil {
		return 0, err
	}
	return lecture.Id, nil
}

func (r *LectureRepo) GetStatistics(totalStudents int, attendanceRepo *AttendanceRepo) (string, error) {
	lectures, err := r.GetAll()
	if err != nil {
		return "", err
	}

	if len(lectures) == 0 {
		return "Лекций ещё не было", nil
	}

	summary := "📊 Статистика посещаемости:\n\n"
	for i, l := range lectures {
		summary += fmt.Sprintf(
			"Лекция %d (%s): %d из %d студентов\n",
			i+1, l.Date, l.CountStudent, totalStudents,
		)
	}

	// самая активная группа
	group, err := attendanceRepo.GetMostActiveGroup()
	if err == nil && group != "" {
		summary += fmt.Sprintf("\n🔥 Самая активная группа: %s", group)
	}

	return summary, nil
}

func (r *LectureRepo) GetLastLectureNumber() (int, error) {
	var lecture models.Lecture

	err := r.db.Order("id DESC").First(&lecture).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil // лекций не было
		}
		return 0, err
	}

	return lecture.Id, nil
}
