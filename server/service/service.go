package service

import (
	"homework_platform/internal/serializer"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	Bind = iota
	BindUri
	None
)

type Service interface {
	Handle(c *gin.Context) (any, error)
}

func HandlerNoBind(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

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

func Handler(s Service) gin.HandlerFunc {
	return HandlerWithBindType(s, Bind)
}

func HandlerBindUri(s Service) gin.HandlerFunc {
	return HandlerWithBindType(s, BindUri)
}

func HandlerWithBindType(s Service, bindType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		
		// Binding using an auto-selected binding engine
		// "application/json" --> JSON binding
		// "application/xml"  --> XML binding
		switch bindType {
		case Bind:
			err = c.ShouldBind(s)
		case BindUri:
			err = c.ShouldBindUri(s)
		case None:
			err = nil
		}
		if err != nil /*&& err != io.EOF*/ {
			log.Printf("[Handler]: Failed to bind: %v\n", err)
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
