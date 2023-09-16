package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	CourseID    int       `json:"course_id"`
	Name        string    `json:"name"`
	TeacherID   string    `json:"teacher_id"`
	BeginDate   time.Time `json:"begin_date"`
	EndDate     time.Time `json:"end_date"`
	Description string    `json:"description"`
}
