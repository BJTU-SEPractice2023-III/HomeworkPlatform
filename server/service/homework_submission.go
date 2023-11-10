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
	HomeworkID uint                    `uri:"id" bind:"required"`
	Content    string                  `form:"content"`
	Files      []*multipart.FileHeader `form:"files"`
}

func (s *SubmitHomework) Handle(c *gin.Context) (any, error) {
	if c.ContentType() != "multipart/form-data" {
		return nil, errors.New("not supported content-type")
	}
	log.Println("ASDSAD")

	var err error
	// 从 Uri 获取 HomeworkID
	err = c.ShouldBindUri(s)
	if err != nil {
		return nil, err
	}
	log.Println("!!")
	log.Println(s)
	// 从 Form 获取其他数据
	err = c.ShouldBind(s)
	if err != nil {
		return nil, err
	}
	log.Println("??")
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
	homworksubmission := models.FindHomeWorkSubmissionByHomeworkIDAndUserID(uint(s.HomeworkID), id.(uint))
	if homworksubmission == nil {
		homworksubmission := models.HomeworkSubmission{
			HomeworkID: uint(s.HomeworkID),
			Content:    s.Content,
			UserID:     id.(uint),
		}
		res := models.AddHomeworkSubmission(&homworksubmission)
		if !res {
			return nil, errors.New("提交失败")
		}
		for _, f := range s.Files {
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
