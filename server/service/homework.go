package service

import (
	"errors"
	"homework_platform/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type AssignHomeworkService struct {
	CourseID    int       `form:"courseid"`
	Name        string    `form:"name"`
	Description string    `form:"description"`
	BeginDate   time.Time `form:"begindate"`
	EndDate     time.Time `form:"enddate"`
}

func (service *AssignHomeworkService) Handle(c *gin.Context) (any, error) {
	_, err := models.CreateHomework(
		service.CourseID,
		service.Name,
		service.Description,
		service.BeginDate,
		service.EndDate,
	)
	print("123")
	if err != nil {
		return nil, errors.New("创建失败")
	}
	file, err := c.FormFile("file") // 根据前端的字段名获取文件
	if err != nil {
		return nil, errors.New("文件上传失败")

	} else {
		filePath := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			return nil, errors.New("文件保存失败")
		}
	}
	return nil, nil
}
