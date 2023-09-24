package service

import (
	"homework_platform/internal/serializer"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Handle(c *gin.Context) (any, error)
}

func Handler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.ShouldBind(s) //检查json和s的结构是否一致
		if err != nil && err != io.EOF {
			c.JSON(http.StatusBadRequest, serializer.ErrorResponse(err))
			return
		}

		res, err := s.Handle(c) //调用实现了接口s的结构体对代码进行处理
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, serializer.ErrorResponse(err))
		} else {
			log.Println("StatusOK")
			c.JSON(http.StatusOK, serializer.Response(res))
		}
	}
}
