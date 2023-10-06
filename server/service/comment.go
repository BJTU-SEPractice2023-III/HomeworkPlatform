package service

import (
	"errors"
	"homework_platform/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentService struct {
	Grade                int    `form:"grade"`
	Comment              string `form:"comment"`
	HomeworkSubmissionID uint   `form:"homeworksubmissionid"`
}

func (service *CommentService) Handle(c *gin.Context) (any, error) {
	if service.Grade < 0 || service.Grade > 100 {
		return nil, errors.New("无效分数")
	}
	homewroksubmission := models.GetHomeWorkSubmissionByID(service.HomeworkSubmissionID)
	homework, res1 := models.GetHomeworkByID(homewroksubmission.HomeworkID)
	if res1 != nil {
		return nil, res1
	}
	if homework.CommentEndDate.Before(time.Now()) {
		return nil, errors.New("超时批阅")
	}
	if homewroksubmission == nil {
		return nil, errors.New("没有找到该作业号")
	}
	id, _ := c.Get("ID")
	comment, res := models.GetCommentByUserIDAndHomeworkSubmissionID(id.(uint), service.HomeworkSubmissionID)
	if res == nil {
		res := comment.(models.Comment).UpdateSelf(service.Comment, service.Grade)
		return nil, res
	}
	return nil, res
}

type GetCommentListsService struct {
	HomeworkID uint `josn:"homeworkid"`
}

func (service *GetCommentListsService) Handle(c *gin.Context) (any, error) {
	id, _ := c.Get("ID")
	err := models.AssignComment(service.HomeworkID)
	if err != nil {
		return nil, err
	}
	commentLists, res := models.GetCommentListsByUserIDAndHomeworknID(id.(uint), service.HomeworkID)
	if res != nil {
		return nil, res
	}

	return commentLists, nil
}
