package models

import (
	"errors"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model

	// A homework submission has many comments
	// Also check homework_submission.go
	// Check: https://gorm.io/docs/has_many.html
	HomeworkSubmissionID uint `json:"homeworkSubmissionId"`
	HomeworkID           uint `json:"homeworkid"`
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
func GetCommentListsByUserIDAndHomeworkID(userid uint, homeworkid uint) (any, error) {
	var comment []Comment
	res := DB.Where("homework_id = ? AND user_id = ?", homeworkid, userid).Find(&comment)
	if res.Error != nil {
		return nil, res.Error
	}
	return comment, nil
}

func CreateComment(HomeworkSubmissionID uint, UserID uint, HomeworkID uint) bool {
	comment := Comment{
		HomeworkSubmissionID: HomeworkSubmissionID,
		UserID:               UserID,
		HomeworkID:           HomeworkID,
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
	if homework.EndDate.Before(time.Now()) {
		//分配作业
		rand.Seed(time.Now().UnixNano()) // 种子随机化
		if homework.Assigned == -1 {
			homework.Assigned = 1
			DB.Save(&homework)
			submissionLists, err := GetSubmissionListsByHomeworkID(HomeworkID)
			if err != nil {
				homework.Assigned = -1
				DB.Save(&homework)
				return err
			}
			//TODO:算法部分,暂时采用每人批三份的方式
			nReviewers := 3 // 每个作业需要三个批改人员
			// var homeworklistsAfterRandon []HomeworkSubmission
			// for _, submission := range submissionLists {
			// 	for i := 0; i < nReviewers; i++ {
			// 		homeworklistsAfterRandon = append(homeworklistsAfterRandon, submission)
			// 	}
			// }
			// for i := len(homeworklistsAfterRandon) - 1; i > 0; i-- {
			// 	//洗牌
			// 	j := rand.Intn(i + 1)
			// 	homeworklistsAfterRandon[i], homeworklistsAfterRandon[j] = homeworklistsAfterRandon[j], homeworklistsAfterRandon[i]
			// }

			// for _, submission := range submissionLists { //在这里获取提交用户的id
			// 	for i := 0; i < nReviewers; i++ {
			// 		CreateComment(homeworklistsAfterRandon[i].ID, submission.UserID, submission.HomeworkID)
			// 	}
			// }
			m := make(map[uint]int)
			var userLists []uint
			for _, submission := range submissionLists {
				m[submission.UserID] = nReviewers
				userLists = append(userLists, submission.UserID)
			}
			for _, submission := range submissionLists { //在这里获取提交用户的id
				for i := 0; i < nReviewers; i++ {
					k := rand.Intn(int(len(userLists)))
					for userLists[k] != submission.UserID && m[userLists[k]] > 0 {
						CreateComment(userLists[k], submission.UserID, submission.HomeworkID)
						m[userLists[k]]--
					}
				}
			}

		}
	} else {
		return errors.New("现在不是批阅时间")
	}
	return nil
}
