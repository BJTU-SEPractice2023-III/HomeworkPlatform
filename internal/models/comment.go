package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model

	// A homework submission has many comments
	// Also check homework_submission.go
	// Check: https://gorm.io/docs/has_many.html
	HomeworkSubmissionID uint `json:"homework_submission_id"`

	// A user has many comments
	// Also check user.go
	// Check: https://gorm.io/docs/has_many.html
	UserID uint `json:"user_id"`

	// Regular fields
	Comment string `json:"comment"`
	Grade   int    `json:"grade"`
}
