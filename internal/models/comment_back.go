package models

import "gorm.io/gorm"

type CommentBack struct {
	gorm.Model
	HomeworkID int `json:"homework_id" gorm:"type:int(20)"`
	CommentID  int `json:"comment_id" gorm:"type:int(20)"`
	Grade      int `json:"grade" gorm:"type:int(20)"`
}
