package models

import (
	"errors"

	"gorm.io/gorm"
)

type Grade struct {
	gorm.Model
	HomeworkSubmissionID uint `json:"submissionid"`
	HomeworkID           uint `json:"homeworkid"`
	Grade                int  `json:"grade"`
}

func GetGradeByID(id uint) (int, error) {
	var grade Grade
	err := DB.Where("homework_submission_id = ?", id).First(&grade).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("grade not found")
		}
		return 0, err
	}
	return grade.Grade, nil
}

func GetGradeListsByHomeworkID(id uint) ([]Grade, error) {
	var grades []Grade
	err := DB.Where("homework_id = ?", id).Find(&grades).Error
	if err != nil {
		return nil, err
	}
	return grades, nil
}
