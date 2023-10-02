package service

import (
	"errors"
	"homework_platform/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type GetGradeBySubmissionIDService struct {
	HomeworkSubmissionID uint `form:"id"`
}

func (service *GetGradeBySubmissionIDService) Handle(c *gin.Context) (any, error) {
	submission := models.GetHomeWorkSubmissionByID(service.HomeworkSubmissionID)
	if submission == nil {
		return nil, errors.New("作业没找到")
	}
	if submission.Final == -1 {
		//分数没有被计算过或者未截止
		grade, res := submission.CalculateGrade()
		if res != nil {
			return nil, res
		}
		homework, res := models.GetHomeworkByID(submission.HomeworkID)
		if res != nil {
			return nil, res
		}
		if homework.EndDate.After(time.Now()) {
			submission.Final = 1
		}
		submission.Grade = grade
		err := submission.UpdateSelf()
		return grade, err
	}
	return submission.Grade, nil
}

type UpdateGradeService struct {
	HomeworkSubmissionID uint `form:"id"`
	Grade                int  `form:"grade"`
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
	submission.Final = 1
	submission.Grade = service.Grade
	err2 := submission.UpdateSelf()
	return nil, err2
}

type GetGradeListsByHomeworkIDService struct {
	HomeworkID uint `form:"id"`
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
	for _, submission := range submissions {
		if submission.Final == -1 {
			//分数没有被计算过或者未截止
			grade, res := submission.CalculateGrade()
			if res != nil {
				return nil, res
			}
			if homework.EndDate.After(time.Now()) {
				submission.Final = 1
			}
			submission.Grade = grade
			err := submission.UpdateSelf()
			if err != nil {
				return nil, err
			}
		}
	}
	return submissions, nil
}