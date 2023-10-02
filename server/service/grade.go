package service

import (
	"homework_platform/internal/models"

	"github.com/gin-gonic/gin"
)

type GetGradeBySubmissionIDService struct {
	HomeworkSubmissionID uint `form:"id"`
}

func (service *GetGradeBySubmissionIDService) Handle(c *gin.Context) (any, error) {
	//TODO:可以设置自己和老师才能查看
	grade, err := models.GetGradeByID(service.HomeworkSubmissionID)
	if err != nil {
		return nil, err
	}
	return grade, nil
}

type GetGradeListsByHomeworkIDService struct {
	HomeworkID uint `form:"id"`
}

func (service *GetGradeListsByHomeworkIDService) Handle(c *gin.Context) (any, error) {
	//TODO:可以检查该课程是不是自己创建的
	// course, err := models.GetCourseByID(int(service.HomeworkID))
	// if err != nil {
	// 	return nil, err
	// }
	// id, _ := c.Get("ID")
	// if course.TeacherID != id {
	// 	return nil, errors.New("不能发布不是您的课程的作业")
	// }
	grades, err := models.GetGradeListsByHomeworkID(service.HomeworkID)
	if err != nil {
		return nil, err
	}
	return grades, nil
}
