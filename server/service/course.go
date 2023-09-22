package service

import (
	"homework_platform/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateCourse struct {
	Name        string    `form:"name"`
	BeginDate   time.Time `form:"begindate"`
	EndDate     time.Time `form:"enddate"`
	Description string    `form:"description"`
	TeacherID   uint      `form:"teacherid"`
}

func (service *CreateCourse) Handle(c *gin.Context) (any, error) {
	err := models.CreateCourse(service.Name, service.BeginDate, service.EndDate, service.Description, int(service.TeacherID))
	if err != nil {
		return nil, err
	}
	return nil, nil
}
