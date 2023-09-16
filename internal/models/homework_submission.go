package models

import "gorm.io/gorm"

type HomeworkSubmission struct {
	gorm.Model
	HomeworkID int    `json:"homework_id" gorm:"type:int(20)"`
	UserID     int    `json:"user_id" gorm:"type:int(20)"`
	Content    string `json:"content"`
	Attachment string `json:"attachment"`
}
