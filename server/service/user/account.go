package user

import (
	"encoding/base64"
	"errors"
	"homework_platform/internal/jwt"
	"homework_platform/internal/models"
	"homework_platform/internal/utils"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Login struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (service *Login) Handle(c *gin.Context) (any, error) {
	// log.Printf("[UserLoginService]: %v, %v\n", service.Username, service.Password)
	var user models.User
	var err error
	if service.Username == "" {
		return nil, errors.New("名称不能为空")
	}
	if service.Password == "" {
		return nil, errors.New("密码不能为空")
	}
	if user, err = models.GetUserByUsername(service.Username); err == gorm.ErrRecordNotFound {
		return nil, err
	}

	if !user.CheckPassword(service.Password) {
		return nil, errors.New("incorrect password")
	}

	var jwtToken string
	jwtToken, err = jwt.CreateToken(user.ID) //根据用id创建jwt
	if err != nil {
		return nil, err
	}

	c.SetCookie("token", jwtToken, 3600, "/", "localhost", false, true)

	res := make(map[string]any)
	res["token"] = jwtToken //之后解码token验证和user是否一致
	res["user"] = user
	// res["user_name"] = user.Username
	// log.Printf("登陆成功")
	return res, nil
}

type UserUpdatePasswordService struct {
	OldPassword string `form:"oldPassword"` // 旧码
	NewPassword string `form:"newPassword"` // 新密码
}

func (service *UserUpdatePasswordService) Handle(c *gin.Context) (any, error) {
	if service.NewPassword == "" {
		return nil, errors.New("密码不能为空")
	}
	id := c.GetUint("ID")
	user, err := models.GetUserByID(id)
	if err != nil {
		return nil, errors.New("该用户不存在")
	}
	// 验证密码
	passwordCheck := user.CheckPassword(service.OldPassword)
	if !passwordCheck {
		return nil, errors.New("密码错误")
	}
	// log.Printf("用户的新密码为%s", service.NewPassword)
	// 修改密码
	if err := user.ChangePassword(service.NewPassword); err != nil {
		return nil, err
	}
	res := make(map[string]any)
	res["msg"] = "修改成功"
	return res, nil
}

// 自己修改密码
type UserselfupdateService struct {
	UserName    string `form:"userName"`
	OldPassword string `form:"oldPassword"` // 旧码
	NewPassword string `form:"newPassword"` // 新密码
}

func (service *UserselfupdateService) Handle(c *gin.Context) (any, error) {
	if service.NewPassword == "" {
		return nil, errors.New("密码不能为空")
	}
	user, err := models.GetUserByUsername(service.UserName)
	if err != nil {
		return nil, errors.New("该用户不存在")
	}
	// 验证密码
	passwordCheck := user.CheckPassword(service.OldPassword)
	if !passwordCheck {
		return nil, errors.New("密码错误")
	}
	// log.Printf("用户的新密码为%s", service.NewPassword)
	// 修改密码
	if err := user.ChangePassword(service.NewPassword); err != nil {
		return nil, err
	}
	res := make(map[string]any)
	res["msg"] = "修改成功"
	return res, nil
}

type GetUserService struct {
	ID uint `uri:"id" binding:"required"`
}

func (service *GetUserService) Handle(c *gin.Context) (any, error) {
	return models.GetUserByID(service.ID)
}

type GetUserNameService struct {
}

func (service *GetUserNameService) Handle(c *gin.Context) (any, error) {
	id := c.GetUint("ID")
	user, err := models.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user.Username, nil
}

type GetAvatar struct {
	ID uint `uri:"id" binding:"required"`
}

func (service *GetAvatar) Handle(c *gin.Context) (any, error) {
	user, err := models.GetUserByID(service.ID)
	if err != nil {
		return nil, err
	}
	return user.Avatar, err
}

type ChangeAvatar struct {
	Avatar *multipart.FileHeader `form:"avatar"`
}

func (s *ChangeAvatar) Handle(c *gin.Context) (any, error) {
	if c.ContentType() != "multipart/form-data" {
		return nil, errors.New("not supported content-type")
	}

	// 从 Form 获取其他数据
	err := c.ShouldBind(s) //获得图片
	if err != nil {
		return nil, err
	}
	if s.Avatar.Size > 1100000 {
		return nil, errors.New("上传图片不可超过1mb")
	}
	// 判断是不是图片
	extension := filepath.Ext(s.Avatar.Filename)
	if !strings.Contains(".jpg.jpeg.png.gif.bmp", extension) {
		return nil, errors.New("unsupported file type")
	}
	id := c.GetUint("ID")
	user, err := models.GetUserByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	avatarByte, err := utils.FileHeaderToBytes(s.Avatar)
	if err != nil {
		return "", err
	}
	base64str := base64.StdEncoding.EncodeToString(avatarByte)
	if err != nil {
		return nil, err
	}
	err = user.ChangeAvatar(base64str)
	return nil, err
}

type Register struct {
	Username string `form:"username"` // 用户名
	Password string `form:"password"` // 密码
}

func (service *Register) Handle(c *gin.Context) (any, error) {
	if len(service.Username) == 0 {
		return nil, errors.New("用户名不能为空")
	}
	if len(service.Password) == 0 {
		return nil, errors.New("密码不能为空")
	}
	_, err := models.CreateUser(service.Username, service.Password)
	return nil, err
}

type GetUserCourses struct {
	ID uint `uri:"id" binding:"required"`
}

func (service *GetUserCourses) Handle(c *gin.Context) (any, error) {
	user, err := models.GetUserByID(service.ID)
	if err != nil {
		return nil, err
	}
	return user.GetCourses()
}

type UpdateSignature struct {
	Signature string `form:"signature"`
}

func (Service *UpdateSignature) Handle(c *gin.Context) (any, error) {
	id, exist := c.Get("ID")
	if !exist {
		return nil, errors.New("不存在id")
	}
	user, err := models.GetUserByID(id.(uint))
	if err != nil {
		return nil, err
	}
	if err := user.ChangeSignature(Service.Signature); err != nil {
		return nil, err
	}
	return nil, nil
}
