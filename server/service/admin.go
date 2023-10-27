package service

import (
	"errors"
	"homework_platform/internal/models"

	// "homework_platform/internal/utils"

	"github.com/gin-gonic/gin"
)

type GetUsersService struct{}

func (service *GetUsersService) Handle(c *gin.Context) (any, error) {
	return models.GetUserList()
}

type DeleteUserService struct {
	ID uint `uri:"id" binding:"required"`
}

func (service *DeleteUserService) Handle(c *gin.Context) (any, error) {
	user, err := models.GetUserByID(service.ID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	result := user.DeleteSelf()
	if !result {
		return nil, errors.New("删除失败")
	}
	res := make(map[string]any)
	res["msg"] = "删除成功"
	return res, nil
}

type UserUpdateService struct { //管理员修改密码
	Username string `form:"username"` // 用户名
	Password string `form:"password"` //新密码
}

func (service *UserUpdateService) Handle(c *gin.Context) (any, error) {
	user, err := models.GetUserByUsername(service.Username)
	if err != nil {
		return nil, errors.New("该用户不存在")
	}
	result := user.ChangePassword(service.Password)
	if !result {
		return nil, errors.New("修改失败")
	}
	res := make(map[string]any)
	res["msg"] = "修改成功"
	return res, nil
}
