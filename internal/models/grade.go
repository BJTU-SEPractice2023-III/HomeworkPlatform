package models

import "gorm.io/gorm"

type Grade struct {
	gorm.Model
	HomeworkSubmissionID int `json:"homework_submission_id" gorm:"type:int(20)"`
	Grade                int `json:"grade" gorm:"type:int(20)"`
}
