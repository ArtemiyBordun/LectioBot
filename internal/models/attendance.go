package models

type Attendance struct {
	Id        int `gorm:"primaryKey"`
	ChatId    int64
	LectureId int
	Number    int
	LastName  string
	Group     string
}

func CreateAttendance(lectureId, number int, chatId int64, lastName, group string) *Attendance {
	return &Attendance{
		ChatId:    chatId,
		LectureId: lectureId,
		Number:    number,
		LastName:  lastName,
		Group:     group,
	}
}
