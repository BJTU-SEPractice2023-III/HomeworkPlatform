package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Homework struct {
	gorm.Model
	CourseID    int       `json:"course_id" gorm:"type:int(20)"`
	Name        string    `json:"name" gorm:"type:varchar(255)"`
	Description string    `json:"description"`
	BeginDate   time.Time `json:"begin_date"`
	EndDate     time.Time `json:"end_date"`

	// A homework has many submissions
	// Also check homeworkSubmission.go
	// Check: https://gorm.io/docs/has_many.html
	HomeworkSubmissions []HomeworkSubmission `json:"-"`
}

func CreateHomework(id int, name string, description string,
	begindate time.Time, endtime time.Time) (any, error) {
	newhomework := Homework{
		CourseID:    id,
		Name:        name,
		Description: description,
		BeginDate:   begindate,
		EndDate:     endtime,
	}
	res := DB.Create(&newhomework)
	if res.Error != nil {
		return nil, errors.New("创建失败")
	}

	return nil, nil
}
