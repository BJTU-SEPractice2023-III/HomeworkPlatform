package service

import (
	"homework_platform/internal/serializer"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	None = 0
	Bind = 1 << iota
	BindUri
)

type Service interface {
	Handle(c *gin.Context) (any, error)
}

func HandlerNoBind(s Service) gin.HandlerFunc {
	return HandlerWithBindType(s, None)
}

func HandlerBind(s Service) gin.HandlerFunc {
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
		if bindType&BindUri != 0 {
			if err = c.ShouldBindUri(s); err != nil {
				// log.Printf("[Handler]: Failed to bind: %v\n", err)
				c.JSON(http.StatusBadRequest, serializer.ErrorResponse(err))
			}
		}
		if bindType&Bind != 0 {
			if err = c.ShouldBind(s); err != nil {
				// log.Printf("[Handler]: Failed to bind: %v\n", err)
				c.JSON(http.StatusBadRequest, serializer.ErrorResponse(err))
			}
		}

		res, err := s.Handle(c)
		if err != nil {
			// log.Println(err.Error())
			c.JSON(http.StatusBadRequest, serializer.ErrorResponse(err))
		} else {
			// log.Println("StatusOK")
			c.JSON(http.StatusOK, serializer.Response(res))
		}
	}
}
