package models

import (
	"fmt"
	"homework_platform/internal/bootstrap"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type HomeworkSubmission struct {
	gorm.Model

	// A homework has many homework submission
	// Also check homework.go
	// Check: https://gorm.io/docs/has_many.html
	HomeworkID uint `json:"homeworkId"`

	// A User has many homework submission
	// Also check user.go
	// Check: https://gorm.io/docs/has_many.html
	UserID    uint     `json:"userId"`
	FilePaths []string `json:"file_paths" gorm:"-"`
	// Regular fields
	Content        string      `json:"content"`
	Score          int         `json:"score" gorm:"default:-1"` //-1表示不是最终结果
	Comments       []Comment   `josn:"-" gorm:"constraint:OnDelete:CASCADE"`
	Complaints []Complaint `josn:"-" gorm:"constraint:OnDelete:CASCADE"`
}

func (homeworksubmission *HomeworkSubmission) GetFiles() {
	root := fmt.Sprintf("./data/homeworkassign/%d/", homeworksubmission.ID)
	files, err := ioutil.ReadDir(root)
	if err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			homeworksubmission.FilePaths = append(homeworksubmission.FilePaths, filepath.Join(root, file.Name()))
		}
	}
}

func GetHomeWorkSubmissionListsByHomeworkID(homeworkID uint) []HomeworkSubmission {
	var homework_submission []HomeworkSubmission
	res := DB.Where("homework_id = ?", homeworkID).Find(&homework_submission)
	if res.Error != nil {
		return nil
	}
	return homework_submission
}

// TODO:后续测试,计算成绩
func (submission *HomeworkSubmission) CalculateGrade() {
	//查询到所有的comment
	homewrork, err := GetHomeworkByID(submission.HomeworkID)
	if homewrork.CommentEndDate.Before(time.Now()) {
		return
	}
	if err != nil {
		return
	}
	comments, res := GetCommentBySubmissionID(submission.ID)
	if res != nil {
		return
	}
	if len(comments) == 0 {
		return
	}
	grade := 0.0
	totalDegree := 0.0
	totalDegreeWithoutDegree := 0
	var userList []User
	var gradeList []int
	for _, comment := range comments {
		if comment.Score == -1 {
			continue
		}
		user, err := GetUserByID(comment.UserID)
		if err != nil {
			return
		}
		totalDegree += user.DegreeOfConfidence
		totalDegreeWithoutDegree += comment.Score
		userList = append(userList, user)
		gradeList = append(gradeList, comment.Score)
		grade += float64(comment.Score) * float64(user.DegreeOfConfidence) //TODO:在这里进行算法开发
	}
	if grade == 0 {
		return
	}
	flag := false
	if submission.Score == -1 {
		flag = true
	}
	average := float64(grade) / totalDegree
	average = math.Round(average)
	submission.UpdateGrade(int(average))
	if flag {
		//TODO:更新degree,不过我是懒B,只算第一次计算
		for i := 0; i < len(userList); i++ {
			userList[i].UpdateDegree(int(average), gradeList[i])
		}
	}

}

// TODO:后续测试,计算成绩
func (submission *HomeworkSubmission) UpdateGrade(Score int) error {
	submission.Score = Score
	return DB.Save(&submission).Error
}

func (homeworksubmission HomeworkSubmission) UpdateSelf() error {
	log.Printf("正在修改homeoworksubmission<id:%d>", homeworksubmission.ID)
	return DB.Save(&homeworksubmission).Error
}

func AddHomeworkSubmission(work *HomeworkSubmission) bool {
	log.Printf("正在创建homeworksubmission<user_id:%d,homework_id:%d>", work.UserID, work.HomeworkID)
	if bootstrap.Sqlite {
		_, err := GetUserByID(work.UserID)
		if err != nil {
			log.Printf("homeworksubmission<user_id:%d,homework_id:%d>:user not exist!", work.UserID, work.HomeworkID)
			return false
		}
		_, err = GetHomeworkByID(work.HomeworkID)
		if err != nil {
			log.Printf("homeworksubmission<user_id:%d,homework_id:%d>:homework not exist!", work.UserID, work.HomeworkID)
			return false
		}
	}
	res := DB.Create(&work)
	return res.Error == nil
}

func GetHomeWorkSubmissionByHomeworkIDAndUserID(homeworkID uint, userID uint) *HomeworkSubmission {
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
	homewroksubmission.GetFiles()
	return &homewroksubmission
}

func GetSubmissionListsByHomeworkID(id uint) ([]HomeworkSubmission, error) {
	var submission []HomeworkSubmission
	if err := DB.Where("homework_id = ?", id).Find(&submission).Error; err != nil {
		return nil, err
	}
	for i := 0; i < len(submission); i++ {
		submission[i].GetFiles()
	}
	return submission, nil
}
