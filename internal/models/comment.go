package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	HomeworkSubmissionID int    `json:"homework_submission_id" gorm:"type:int(20)"`
	Comment              string `json:"comment"`
	ReviewerID           int    `json:"reviewer_id" gorm:"type:int(20)"`
	Grade                int    `json:"grade"`
}
