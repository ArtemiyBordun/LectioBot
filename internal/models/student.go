package models

import "time"

type Student struct {
	ChatId    int64  `gorm:"primaryKey"`
	Name      string //ФИО
	UserName  string //в тг
	Group     string
	CreatedAt time.Time
}

func CreateStudent(chatId int64, name, userName, group string) *Student {
	return &Student{
		ChatId:    chatId,
		Name:      name,
		UserName:  userName,
		Group:     group,
		CreatedAt: time.Now(),
	}
}
