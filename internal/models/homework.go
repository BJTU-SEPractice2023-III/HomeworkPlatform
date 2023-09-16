package models

import (
	"time"

	"gorm.io/gorm"
)

type Homework struct {
	gorm.Model
	CourseID    int       `json:"course_id" gorm:"type:int(20)"`
	Name        string    `json:"name" gorm:"type:varchar(255)"`
	Description string    `json:"description" gorm:"type:varchar(255)"`
	BeginDate   time.Time `json:"begin_date"`
	EndDate     time.Time `json:"end_date"`
}
