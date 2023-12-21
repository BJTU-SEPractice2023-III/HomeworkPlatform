package service

import (
	"errors"
	"homework_platform/internal/models"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Download file by id
type DownloadFileById struct {
	FileId uint `uri:"id" binding:"required"`
}

func (service *DownloadFileById) Handle(c *gin.Context) (any, error) {
	c.Header("Content-Disposition", "attachment")

	file, err := models.GetFileByID(service.FileId)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(file.Path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 将文件内容写入HTTP响应
	_, err = io.Copy(c.Writer, f)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type GetFileService struct {
	// FilePath string `uri:"path" binding:"required"`
}

func (service *GetFileService) Handle(c *gin.Context) (any, error) {
	c.Header("Content-Disposition", "attachment")
	filePath := c.Param("path")[1:]
	println(filePath)
	if !strings.HasPrefix(filePath, "./data") && !strings.HasPrefix(filePath, "data") {
		return nil, errors.New("无法访问该文件")
	}
	log.Printf("path:%s", filePath)
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 将文件内容写入HTTP响应
	_, err = io.Copy(c.Writer, f)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
