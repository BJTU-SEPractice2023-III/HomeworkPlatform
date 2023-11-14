package service

import (
	"errors"
	"homework_platform/internal/models"

	"github.com/gin-gonic/gin"
)

type GetGradeBySubmissionIDService struct {
	HomeworkSubmissionID uint `uri:"id" binding:"required"`
}

func (service *GetGradeBySubmissionIDService) Handle(c *gin.Context) (any, error) {
	submission := models.GetHomeWorkSubmissionByID(service.HomeworkSubmissionID)
	if submission == nil {
		return nil, errors.New("作业没找到")
	}
	return submission.Score, nil
}

type UpdateGradeService struct {
	HomeworkSubmissionID uint `form:"homeworksubmissionid"`
	Score                int  `form:"score"`
}

func (service *UpdateGradeService) Handle(c *gin.Context) (any, error) {
	submission := models.GetHomeWorkSubmissionByID(service.HomeworkSubmissionID)
	if submission == nil {
		return nil, errors.New("作业没找到")
	}
	homework, res := models.GetHomeworkByID(submission.HomeworkID)
	if res != nil {
		return nil, res
	}
	course, err := models.GetCourseByID(homework.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if id.(uint) != course.TeacherID {
		return nil, errors.New("您无权限修改")
	}
	submission.Score = service.Score
	// submission.Grade = service.Final
	err2 := submission.UpdateSelf()
	//TODO:还需要更新那些人的置信度,但是我是懒B
	return nil, err2
}

type GetGradeListsByHomeworkIDService struct {
	HomeworkID uint `uri:"id" binding:"required"`
}

type MyMap struct {
	UserID   uint   `form:"userid"`
	UserName string `form:"username"`
	Score    int    `form:"Score"`
}

func (service *GetGradeListsByHomeworkIDService) Handle(c *gin.Context) (any, error) {
	homework, res := models.GetHomeworkByID(service.HomeworkID)
	if res != nil {
		return nil, res
	}
	course, err := models.GetCourseByID(homework.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if id.(uint) != course.TeacherID {
		return nil, errors.New("您无权限查询")
	}
	submissions, err2 := models.GetSubmissionListsByHomeworkID(service.HomeworkID)
	if err2 != nil {
		return nil, err2
	}
	var maps []MyMap
	for _, submission := range submissions {
		user, err := models.GetUserByID(submission.UserID)
		if err != nil {
			return nil, err
		}
		maps = append(maps, MyMap{UserID: user.ID, UserName: user.Username, Score: submission.Score})
	}
	return maps, nil
}
