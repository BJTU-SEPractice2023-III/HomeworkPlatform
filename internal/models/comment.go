package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model

	// A homework submission has many comments
	// Also check homework_submission.go
	// Check: https://gorm.io/docs/has_many.html
	HomeworkSubmissionID uint `json:"homeworkSubmissionId"`

	// A user has many comments
	// Also check user.go
	// Check: https://gorm.io/docs/has_many.html
	UserID uint `json:"userId"`

	// Regular fields
	Comment string `json:"comment"`
	Grade   int    `json:"grade"`
}

func (comment Comment) UpdateSelf(comm string, grade int) error {
	res := DB.Model(&comment).Updates(Comment{Comment: comm, Grade: grade})
	return res.Error
}

func GetCommentBySubmissionID(submissionid uint) ([]Comment, error) {
	var comments []Comment
	res := DB.Where("homework_submission_id = ?", submissionid).Find(&comments)
	if res.Error != nil {
		return nil, res.Error
	}
	return comments, nil
}

func GetCommentByUserIDAndHomeworkSubmissionID(userid uint, homeworksubmissionid uint) (any, error) {
	var comment Comment
	res := DB.Where("homework_submission_id = ? AND user_id = ?", homeworksubmissionid, userid).First(&comment)
	if res.Error != nil {
		return nil, res.Error
	}
	return comment, nil
}

func CreateComment(HomeworkSubmissionID uint, UserID uint, Commen string, Grade int) bool {
	comment := Comment{
		HomeworkSubmissionID: HomeworkSubmissionID,
		UserID:               UserID,
		Comment:              Commen,
		Grade:                Grade,
	}
	res := DB.Create(&comment)
	return res.Error == nil
}

func AssignComment(HomeworkID uint) error {
	//在这里我们进行作业的分配,每次如果作业没有被分配并且时间到了那么我们就分配!
	homework, err := GetHomeworkByID(HomeworkID)
	if err != nil {
		return err
	}
	if homework.EndDate.Before(time.Now()) && homework.Assigned == -1 {
		//分配作业
		homework.Assigned = 1
		DB.Save(&homework)

	}
	return nil
}
