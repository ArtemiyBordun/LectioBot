package models

type Lecture struct {
	Id           int `gorm:"primaryKey"`
	CountStudent int
	Date         string
}

func CreateLecture(id, countStudent int, date string) *Lecture {
	return &Lecture{
		Id:           id,
		CountStudent: countStudent,
		Date:         date,
	}
}
