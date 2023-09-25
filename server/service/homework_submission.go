package service

import (
	"errors"
	"fmt"
	"homework_platform/internal/models"
	"log"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
)

type SubmitHomework struct {
	HomeworkID int                     `form:"homeworkid"`
	Content    string                  `form:"content"`
	SubmitTime time.Time               `form:"submittime"`
	Files      []*multipart.FileHeader `form:"files"`
}

func (service *SubmitHomework) Handle(c *gin.Context) (any, error) {
	id, _ := c.Get("ID")
	homework, err2 := models.GetHomeworkByID(uint(service.HomeworkID))
	if err2 != nil {
		return nil, err2
	}
	course, err := models.GetCourseByID(homework.CourseID)
	if err != nil {
		return nil, err
	}
	if !course.FindStudents(id.(uint)) {
		return nil, errors.New("请先选择这门课")
	}
	homework_submission := models.HomeworkSubmission{
		HomeworkID: uint(service.HomeworkID),
		Content:    service.Content,
		UserID:     id.(uint),
	}
	res := models.AddHomeworkSubmission(&homework_submission)
	if !res {
		return nil, errors.New("提交失败")
	}
	for _, f := range service.Files {
		log.Println(f.Filename)
		dst := fmt.Sprintf("./data/homework_submission/%d//%s", homework_submission.ID, f.Filename)
		// 上传文件到指定的目录
		c.SaveUploadedFile(f, dst)
	}
	return homework_submission, nil
}
