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
	if service.Score < 0 || service.Score > 100 {
		return nil, errors.New("无效分数")
	}
	if service.Comment == "" {
		return nil, errors.New("评论为空")
	}
	homewroksubmission, _ := models.GetHomeworkSubmissionById(service.HomeworkSubmissionID)
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
	comment, err := models.GetCommentByUserIDAndHomeworkSubmissionID(id.(uint), service.HomeworkSubmissionID)
	if err == nil {
		res := comment.UpdateSelf(service.Comment, service.Score)
		num, len := models.GetCommentNum(service.HomeworkSubmissionID)
		if num == len {
			homewroksubmission.CalculateGrade()
		}
		return nil, res
	}
	return nil, err
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
	if homework.EndDate.After(time.Now()) {
		return nil, errors.New("评阅未开始")
	}

	id := c.GetUint("ID")
	if id == course.TeacherID {
		commentList, err := models.GetCommentsByHomeworkId(service.HomeworkID)
		if err != nil {
			return nil, err
		}
		return commentList, nil
	}

	commentList, err := homework.GetCommentsByUserId(id)
	if err != nil {
		return nil, err
	}
	var homeworkSubmissions []models.HomeworkSubmission
	for _, comment := range commentList {
		homeworkSubmission, err := models.GetHomeworkSubmissionById(comment.HomeworkSubmissionID)
		if err == nil {
			homeworkSubmissions = append(homeworkSubmissions, *homeworkSubmission)
		}
	}
	m := make(map[string]any)
	m["homework_submission"] = homeworkSubmissions
	m["comment_lists"] = commentList
	log.Printf("%x", len(commentList))
	return m, nil
}

type GetMyCommentService struct {
	HomeworkID uint `uri:"id" binding:"required"`
}

func (service *GetMyCommentService) Handle(c *gin.Context) (any, error) {
	homework, err := models.GetHomeworkByID(service.HomeworkID)
	if err != nil {
		return nil, err
	}
	if homework.EndDate.After(time.Now()) {
		return nil, errors.New("评阅未开始")
	}
	id := c.GetUint("ID")
	submission, err := homework.GetSubmissionByUserId(id)
	if err != nil {
		return nil, err
	}
	comments, err := models.GetCommentBySubmissionID(submission.ID)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
