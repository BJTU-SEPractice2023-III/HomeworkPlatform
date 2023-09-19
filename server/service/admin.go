package service

import (
	"homework_platform/internal/models"
	// "homework_platform/internal/utils"

	"github.com/gin-gonic/gin"
)

type GetUsersService struct{}

func (service *GetUsersService) Handle(c *gin.Context) (any, error) {
	return models.GetUserList()
}
