package service

import (
	"errors"
	"homework_platform/internal/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentService struct {
	Score                int    `form:"score"`
	Comment              string `form:"comment"`
	HomeworkSubmissionID uint   `uri:"id" binding:"required"`
}

func (service *CommentService) Handle(c *gin.Context) (any, error) {
	err := c.ShouldBindUri(service)
	if err != nil {
		return nil, err
	}
	//绑定reason
	err = c.ShouldBind(service)
	if err != nil {
		return nil, err
	}
	if service.Score < 0 || service.Score > 100 {
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
	// comment是预先分配好的,所以不需要自我创建
	comment, res := models.GetCommentByUserIDAndHomeworkSubmissionID(id.(uint), service.HomeworkSubmissionID)
	if res == nil {
		res := comment.(models.Comment).UpdateSelf(service.Comment, service.Score)
		num := models.GetCommentNum(service.HomeworkSubmissionID)
		if num == 3 {
			homewroksubmission.CalculateGrade()
		}
		return nil, res
	}
	return nil, res
}

type GetCommentListsService struct {
	HomeworkID uint `uri:"id" binding:"required"`
}

func (service *GetCommentListsService) Handle(c *gin.Context) (any, error) {
	log.Println("[GetCommentListsService]: Trying to assign comments")
	if err := models.AssignComment(service.HomeworkID); err != nil {
		return nil, err
	}

	homework, err := models.GetHomeworkByID(service.HomeworkID)
	if err != nil {
		return nil, err
	}
	course, err := models.GetCourseByID(homework.CourseID)
	if err != nil {
		return nil, err
	}

	id, _ := c.Get("ID")
	if id == course.TeacherID {
		commentList, err := models.GetCommentsByHomeworkId(service.HomeworkID)
		if err != nil {
			return nil, err
		}
		return commentList, nil
	}

	commentList, err := models.GetCommentListsByUserIDAndHomeworkID(id.(uint), service.HomeworkID)
	if err != nil {
		return nil, err
	}
	var homework_submission []models.HomeworkSubmission
	for _, comment := range commentList {
		homework_submission = append(homework_submission, *models.GetHomeWorkSubmissionByID(comment.HomeworkSubmissionID))
	}
	m := make(map[string]any)
	m["homework_submission"] = homework_submission
	m["comment_lists"] = commentList
	log.Printf("%x", len(commentList))
	return m, nil
}

type GetMyCommentService struct {
	HomeworkID uint `uri:"id" binding:"required"`
}

func (service *GetMyCommentService) Handle(c *gin.Context) (any, error) {
	_, err := models.GetHomeworkByID(service.HomeworkID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	submission := models.GetHomeWorkSubmissionByHomeworkIDAndUserID(service.HomeworkID, id.(uint))
	comments, err := models.GetCommentBySubmissionID(submission.ID)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
