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

	var err error
	// 从 Uri 获取 HomeworkID
	err = c.ShouldBindUri(s)
	if err != nil {
		return nil, err
	}
	// 从 Form 获取其他数据
	err = c.ShouldBind(s)
	if err != nil {
		return nil, err
	}

	id, _ := c.Get("ID")
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
	log.Printf("正在检查该用户有没有选课%d", id.(uint))
	if !course.FindStudents(id.(uint)) {
		return nil, errors.New("请先选择这门课")
	}
	homworksubmission := models.GetHomeWorkSubmissionByHomeworkIDAndUserID(uint(s.HomeworkID), id.(uint))
	if homworksubmission == nil {
		homworksubmission := models.HomeworkSubmission{
			HomeworkID: uint(s.HomeworkID),
			Content:    s.Content,
			UserID:     id.(uint),
		}
		// res := models.AddHomeworkSubmission(&homworksubmission)
		_, err := homework.AddSubmission(homworksubmission)
		if err != nil {
			return nil, err
		}
		for _, f := range s.Files {
			log.Println(f.Filename)
			dst := fmt.Sprintf("./data/homework_submission/%d/%s", homworksubmission.ID, f.Filename)
			// 上传文件到指定的目录
			c.SaveUploadedFile(f, dst)
		}
		return nil, nil
	}
	return nil, errors.New("不可重复提交")
}

type GetHomeworkSubmission struct {
	HomeworkId uint `uri:"id" binding:"required"`
}

func (service *GetHomeworkSubmission) Handle(c *gin.Context) (any, error) {
	userid, _ := c.Get("ID")
	log.Printf("用户id为%d", userid.(uint))
	log.Printf("homeworkid为%d", service.HomeworkId)
	homework, err := models.GetHomeworkByID(service.HomeworkId)
	if err != nil {
		return "该作业号不存在", nil
	}
	for _, value := range homework.HomeworkSubmissions {
		if value.UserID == userid.(uint) {
			value.GetFiles()
			return value, nil
		}
	}
	return nil, errors.New("该用户未提交作业")
}

type UpdateSubmission struct {
	HomeworkID uint                    `uri:"id" bind:"required"`
	Content    string                  `form:"content"`
	Files      []*multipart.FileHeader `form:"files"`
}

func (s *UpdateSubmission) Handle(c *gin.Context) (any, error) {
	if c.ContentType() != "multipart/form-data" {
		return nil, errors.New("not supported content-type")
	}

	var err error
	// 从 Uri 获取 HomeworkID
	err = c.ShouldBindUri(s)
	if err != nil {
		return nil, err
	}
	// 从 Form 获取其他数据
	err = c.ShouldBind(s)
	if err != nil {
		return nil, err
	}
	log.Println(s)

	id, _ := c.Get("ID")
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
	if !course.FindStudents(id.(uint)) {
		return nil, errors.New("请先选择这门课")
	}
	homworksubmission := models.GetHomeWorkSubmissionByHomeworkIDAndUserID(uint(s.HomeworkID), id.(uint))
	if homworksubmission != nil {
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

type GetSubmissionService struct {
	HomeworkID uint `uri:"id" bind:"required"`
}

func (s *GetSubmissionService) Handle(c *gin.Context) (any, error) {
	submit := models.GetHomeWorkSubmissionByID(s.HomeworkID)
	submit.GetFiles()
	return submit, nil
}
