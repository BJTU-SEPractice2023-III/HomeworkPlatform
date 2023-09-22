package service

import (
	"homework_platform/internal/jwt"
	"homework_platform/internal/models"

	// "homework_platform/internal/utils"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserLoginService struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (service *UserLoginService) Handle(c *gin.Context) (any, error) {
	var user models.User
	var err error

	if user, err = models.GetUserByUsername(service.Username); err == gorm.ErrRecordNotFound {
		return nil, errors.New("not exist")
	}

	if !user.CheckPassword(service.Password) {
		return nil, errors.New("incorrect password")
	}

	var jwtToken string
	jwtToken, err = jwt.CreateToken(user.ID) //根据用id创建jwt

	res := make(map[string]any)
	res["token"] = jwtToken //之后解码token验证和user是否一致
	res["user"] = user
	// res["user_name"] = user.Username

	return res, nil
}

type UserselfUpdateService struct { //自己修改密码
	Username    string `form:"username"`    // 用户名
	OldPassword string `form:"oldpassword"` //旧码
	NewPassword string `form:"newpassword"`
}

func (service *UserselfUpdateService) Handle(c *gin.Context) (any, error) {
	user, err := models.GetUserByUsername(service.Username)
	if err != nil {
		return nil, errors.New("该用户不存在")
	}
	//验证密码
	passwordCheck := user.CheckPassword(service.OldPassword)
	if !passwordCheck {
		return nil, errors.New("密码错误!")
	}
	//修改密码
	result := user.ChangePassword(service.NewPassword)
	if !result {
		return nil, errors.New("修改失败")
	}
	res := make(map[string]any)
	res["msg"] = "修改成功"
	return res, nil
}

type GetUserService struct {
	ID uint `form:"id"`
}

func (service *GetUserService) Handle(c *gin.Context) (any, error) {
	if user, err := models.GetUserByID(service.ID); err == nil {
		return user, nil
	} else {
		return nil, err
	}
}

type UserRegisterService struct {
	Username string `form:"username"` // 用户名
	Password string `form:"password"` // 密码
}

func (service *UserRegisterService) Handle(c *gin.Context) (any, error) {
	res, err := models.CreateUser(service.Username, service.Password)
	if err != nil {
		return nil, err
	}
	return res, nil
}
