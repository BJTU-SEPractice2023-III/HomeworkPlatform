package service

// import (
// 	"errors"
// 	"homework_platform/internal/models"
// 	"log"

// 	// "homework_platform/internal/utils"

// 	"github.com/gin-gonic/gin"
// )

// type GetUsersService struct{}

// func (service *GetUsersService) Handle(c *gin.Context) (any, error) {
// 	return models.GetUsers()
// }

// // DeleteUserService deletes a `User` with `ID`
// type DeleteUserService struct {
// 	ID uint `uri:"id" binding:"required"`
// }

// func (service *DeleteUserService) Handle(c *gin.Context) (res any, err error) {
// 	if err = models.DeleteUserById(service.ID); err != nil {
// 		// log.Printf("删除失败(%s)", err)
// 	}
// 	return
// }

// type UserUpdateService struct { //管理员修改密码
// 	Username string `form:"username"` // 用户名
// 	Password string `form:"password"` //新密码
// }

// func (service *UserUpdateService) Handle(c *gin.Context) (any, error) {
// 	user, err := models.GetUserByUsername(service.Username)
// 	if err != nil {
// 		return nil, errors.New("该用户不存在")
// 	}
// 	if err := user.ChangePassword(service.Password); err != nil {
// 		return nil, errors.New("修改失败")
// 	}
// 	res := make(map[string]any)
// 	res["msg"] = "修改成功"
// 	return res, nil
// }
