package service

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type GetFileService struct {
	FilePath string `form:"filepath"`
}

func (service *GetFileService) Handle(c *gin.Context) (any, error) {
	c.Header("Content-Disposition", "attachment")
	filePath := service.FilePath
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
