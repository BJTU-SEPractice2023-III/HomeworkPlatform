package service

import (
	"errors"
	"fmt"
	"homework_platform/internal/models"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type SubmitHomework struct {
	HomeworkID uint                    `uri:"id" bind:"required"`
	Content    string                  `form:"content"`
	Files      []*multipart.FileHeader `form:"files"`
}

func (s *SubmitHomework) Handle(c *gin.Context) (any, error) {
	if c.ContentType() != "multipart/form-data" {
		return nil, errors.New("not supported content-type")
	}

	id := c.GetUint("ID")
	time := time.Now()
	homework, err := models.GetHomeworkByID(uint(s.HomeworkID))
	if err != nil {
		return nil, err
	}
	if homework.EndDate.Before(time) {
		return nil, errors.New("超时提交")
	}
	course, err := models.GetCourseByID(homework.CourseID)
	if err != nil {
		return nil, err
	}
	log.Printf("正在检查该用户有没有选课%d", id)
	if !course.FindStudents(id) {
		return nil, errors.New("请先选择这门课")
	}
	homworksubmission, err := homework.GetSubmissionByUserId(id)
	if homworksubmission == nil || err != nil {
		submission, err := homework.AddSubmission(id, s.Content)
		if err != nil {
			return nil, err
		}
		for _, f := range s.Files {
			file, err := models.CreateFileFromFileHeaderAndContext(f, c)
			if err != nil {
				// TODO: err handle
			} else {
				file.Attach(submission.ID, models.TargetTypeHomeworkSubmission)
			}
		}
		return nil, nil
	}
	return nil, errors.New("不可重复提交")
}

type GetHomeworkUserSubmission struct {
	HomeworkId uint `uri:"id" binding:"required"`
}

func (service *GetHomeworkUserSubmission) Handle(c *gin.Context) (any, error) {
	userId := c.GetUint("ID")
	// log.Printf("用户id为%d", userId)
	// log.Printf("homeworkid为%d", service.HomeworkId)
	homework, err := models.GetHomeworkByID(service.HomeworkId)
	if err != nil {
		return "该作业号不存在", nil
	}
	submission, err := homework.GetSubmissionByUserId(userId)
	if err != nil {
		return nil, err
	}
	return *submission, nil
}

type UpdateUserSubmission struct {
	HomeworkID uint                    `uri:"id" bind:"required"`
	Content    string                  `form:"content"`
	Files      []*multipart.FileHeader `form:"files"`
}

func (s *UpdateUserSubmission) Handle(c *gin.Context) (any, error) {
	if c.ContentType() != "multipart/form-data" {
		return nil, errors.New("not supported content-type")
	}

	id := c.GetUint("ID")
	time := time.Now()
	homework, err := models.GetHomeworkByID(uint(s.HomeworkID))
	if err != nil {
		return nil, err
	}
	if homework.EndDate.Before(time) {
		return nil, errors.New("超时提交")
	}
	course, err := models.GetCourseByID(homework.CourseID)
	if err != nil {
		return nil, err
	}
	if !course.FindStudents(id) {
		return nil, errors.New("请先选择这门课")
	}
	homworksubmission, err := homework.GetSubmissionByUserId(id)
	if homworksubmission != nil || err != nil {
		homworksubmission.Content = s.Content
		homworksubmission.UpdateSelf()
		os.RemoveAll(fmt.Sprintf("./data/homework_submission/%d", homworksubmission.ID))
		for _, f := range s.Files {
			log.Println(f.Filename)
			dst := fmt.Sprintf("./data/homework_submission/%d/%s", homworksubmission.ID, f.Filename)
			// 上传文件到指定的目录
			c.SaveUploadedFile(f, dst)
		}
		return nil, nil
	}
	return nil, errors.New("请先提交作业")
}

type GetSubmissionByIdService struct {
	HomeworkID uint `uri:"id" bind:"required"`
}

func (s *GetSubmissionByIdService) Handle(c *gin.Context) (any, error) {
	return models.GetHomeworkSubmissionById(s.HomeworkID)
}
