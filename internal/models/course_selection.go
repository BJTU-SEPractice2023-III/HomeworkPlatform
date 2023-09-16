package models

import "gorm.io/gorm"

type CourseSelection struct {
	gorm.Model
	StudentID int `json:"student_id" gorm:"type:int(10)"`
	CourseID  int `json:"course-id" gorm:"type:int(10)"`
}
