package models

import "gorm.io/gorm"

type HomeworkSubmission struct {
	gorm.Model

	// A homework has many homework submission
	// Also check homework.go
	// Check: https://gorm.io/docs/has_many.html
	HomeworkID uint `json:"homework_id"`

	// A User has many homework submission
	// Also check user.go
	// Check: https://gorm.io/docs/has_many.html
	UserID uint `json:"user_id"`

	// Regular fields
	Content string `json:"content"`
}

func AddHomeworkSubmission(work *HomeworkSubmission) bool {
	println(work.UserID)
	res := DB.Create(&work)
	return res.Error == nil
}
