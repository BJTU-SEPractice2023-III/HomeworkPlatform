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
	Files      []*multipart.FileHeader `form:"files"`
}

func (service *SubmitHomework) Handle(c *gin.Context) (any, error) {
	id, _ := c.Get("ID")
	time := time.Now()
	homework, err2 := models.GetHomeworkByID(uint(service.HomeworkID))
	if err2 != nil {
		return nil, err2
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
	homworksubmission := models.FindHomeWorkSubmissionByHomeworkIDAndUserID(uint(service.HomeworkID), id.(uint))
	if homworksubmission == nil {
		homworksubmission := models.HomeworkSubmission{
			HomeworkID: uint(service.HomeworkID),
			Content:    service.Content,
			UserID:     id.(uint),
		}
		res := models.AddHomeworkSubmission(&homworksubmission)
		if !res {
			return nil, errors.New("提交失败")
		}
		for _, f := range service.Files {
			log.Println(f.Filename)
			dst := fmt.Sprintf("./data/homework_submission/%d/%s", homworksubmission.ID, f.Filename)
			// 上传文件到指定的目录
			c.SaveUploadedFile(f, dst)
		}
		return nil, nil
	}
	return nil, nil
}

type GetHomeworkSubmission struct {
	userid     uint `uri:"userid" binding:"required"`
	homeworkid uint `uri:"homeworkid" binding:"required"`
}

func (service *GetHomeworkSubmission) Handle(c *gin.Context) (any, error) {
	homework, err := models.GetHomeworkByIDWithSubmissionLists(service.homeworkid)
	if err != nil {
		return "该作业号不存在", nil
	}
	for _, value := range homework.HomeworkSubmissions {
		if value.UserID == service.userid {
			return value, nil
		}
	}
	return nil, errors.New("该用户未提交作业")
}
