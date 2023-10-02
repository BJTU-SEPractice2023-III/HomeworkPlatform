package models

import (
	"math"

	"gorm.io/gorm"
)

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
	Grade   int    `json:"-" gorm:"default:-1"`
	Final   int    `json:"-" gorm:"default:-1"` //-1表示不是最终结果
}

func (submission *HomeworkSubmission) CalculateGrade() (int, error) {
	//查询到所有的comment
	comments, res := GetCommentBySubmissionID(submission.ID)
	if res != nil {
		return -1, res
	}
	grade := 0
	for _, comment := range comments {
		grade += comment.Grade
	}
	average := float64(grade) / float64(len(comments))
	average = math.Round(average)
	return int(average), nil
}

func (submission *HomeworkSubmission) UpdateGrade(grade int) error {
	submission.Grade = grade
	return DB.Save(&submission).Error
}

func (homeworksubmission HomeworkSubmission) UpdateSelf() error {
	return DB.Save(&homeworksubmission).Error
}

func AddHomeworkSubmission(work *HomeworkSubmission) bool {
	println(work.UserID)
	res := DB.Create(&work)
	return res.Error == nil
}

func FindHomeWorkSubmissionByHomeworkIDAndUserID(homeworkID uint, userID uint) *HomeworkSubmission {
	var submission *HomeworkSubmission
	if err := DB.Where("user_id = ? AND homework_id = ?", userID, homeworkID).First(&submission).Error; err != nil {
		return nil
	}

	if submission.ID != 0 {
		return submission
	} else {
		return nil
	}
}

func GetHomeWorkSubmissionByID(homewroksubmissionid uint) *HomeworkSubmission {
	var homewroksubmission HomeworkSubmission

	res := DB.First(&homewroksubmission, homewroksubmissionid)
	if res.Error != nil {
		return nil
	}
	return &homewroksubmission
}

func GetSubmissionListsByHomeworkID(id uint) ([]HomeworkSubmission, error) {
	var submission []HomeworkSubmission
	if err := DB.Where("homework_id = ?", id).First(&submission).Error; err != nil {
		return submission, err
	}
	return submission, nil
}
