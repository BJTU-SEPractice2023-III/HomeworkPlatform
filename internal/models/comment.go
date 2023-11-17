package models

import (
	"errors"
	"homework_platform/internal/bootstrap"
	"log"
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
	Score   int    `json:"score" gorm:"default:-1"`
}

// 计算已经有几个人批过了
func GetCommentNum(homeworksubmission_id uint) uint {
	comments, err := GetCommentBySubmissionID(homeworksubmission_id)
	if err != nil {
		return 0
	}
	var num uint
	num = 0
	for _, comment := range comments {
		if comment.Score != -1 {
			num += 1
		}
	}
	return num
}

func (comment Comment) UpdateSelf(comm string, score int) error {
	res := DB.Model(&comment).Updates(Comment{Comment: comm, Score: score})
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

func GetCommentListsByUserIDAndHomeworkID(userid uint, homeworkid uint) ([]Comment, error) {
	var comment []Comment
	res := DB.Where("homework_id = ? AND user_id = ?", homeworkid, userid).Find(&comment)
	if res.Error != nil {
		return nil, res.Error
	}
	return comment, nil
}

func CreateComment(HomeworkSubmissionID uint, UserID uint, HomeworkID uint) bool {
	log.Printf("正在创建comment<user_id:%d,homework_submission_id:%d>", UserID, HomeworkSubmissionID)
	if bootstrap.Sqlite {
		_, err := GetUserByID(UserID)
		if err != nil {
			log.Printf("用户<user_id:%d>不存在", UserID)
			return false
		}
		res := GetHomeWorkSubmissionByID(HomeworkSubmissionID)
		if res == nil {
			log.Printf("作业提交<submission:id:%d>不存在", HomeworkSubmissionID)
			return false
		}
	}

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
			homework.Assigned = 1 //标志位,表示是否已经被分配
			DB.Save(&homework)
			submissionLists, err := GetSubmissionListsByHomeworkID(HomeworkID)
			if err != nil {
				homework.Assigned = -1
				DB.Save(&homework)
				return err
			}
			//TODO:算法部分,暂时采用每人批三份的方式
			nReviewers := 3 // 每个作业需要三个批改人员
			m := make(map[uint]int)
			var userLists []uint
			if len(submissionLists) <= nReviewers {
				//少于3人,那么直接分配其他的人就行
				for _, users := range submissionLists {
					for _, submission := range submissionLists {
						if users.ID != submission.ID {
							CreateComment(submission.ID, users.UserID, submission.HomeworkID)
						}
					}
				}
			} else {
				for _, submission := range submissionLists {
					m[submission.UserID] = nReviewers
					userLists = append(userLists, submission.UserID)
				}
				for _, submission := range submissionLists { //在这里获取提交用户的id
					var used []uint
					for i := 0; i < nReviewers; i++ {
						for {
							k := rand.Intn(int(len(userLists)))
							found := false
							for _, z := range used {
								if int(z) == k {
									found = true
									break
								}
							}
							if userLists[k] != submission.UserID && m[userLists[k]] > 0 && !found {
								CreateComment(submission.ID, userLists[k], submission.HomeworkID)
								used = append(used, uint(k))
								m[userLists[k]]--
								break
							}
						}
					}
				}
			}
		}
	} else {
		return errors.New("现在不是批阅时间")
	}
	return nil
}
