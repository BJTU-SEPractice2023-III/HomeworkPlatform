package service

import (
	"errors"
	"homework_platform/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateCourse struct {
	Name        string    `form:"name"`
	BeginDate   time.Time `form:"begindate"`
	EndDate     time.Time `form:"enddate"`
	Description string    `form:"description"`
}

func (service *CreateCourse) Handle(c *gin.Context) (any, error) {
	id, _ := c.Get("ID")
	err := models.CreateCourse(service.Name, service.BeginDate, service.EndDate, service.Description, id.(uint))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type UpdateCourseDescription struct {
	CourseID    int    `form:"courseid"`
	Description string `form:"description"`
}

func (service *UpdateCourseDescription) Handle(c *gin.Context) (any, error) {
	println(service.Description)
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能修改不是您的课程")
	}
	res := course.UpdateCourseDescription(service.Description)
	if res == false {
		return nil, errors.New("创建失败")
	}
	return nil, nil
}
