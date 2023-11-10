package service

import (
	"homework_platform/internal/jwt"
	"homework_platform/internal/models"
	"time"

	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserLoginService struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (service *UserLoginService) Handle(c *gin.Context) (any, error) {
	log.Printf("[UserLoginService]: %v, %v\n", service.Username, service.Password)
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
	if err != nil {
		return nil, err
	}

	res := make(map[string]any)
	res["token"] = jwtToken //之后解码token验证和user是否一致
	res["user"] = user
	// res["user_name"] = user.Username

	return res, nil
}

// 自己修改密码
type UserselfUpdateService struct {
	Username    string `form:"username"`    // 用户名
	OldPassword string `form:"oldpassword"` // 旧码
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
		return nil, errors.New("密码错误")
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
	ID uint `uri:"id" binding:"required"`
}

func (service *GetUserService) Handle(c *gin.Context) (any, error) {
	return models.GetUserByID(service.ID)
}

type UserRegisterService struct {
	Username string `form:"username"` // 用户名
	Password string `form:"password"` // 密码
}

func (service *UserRegisterService) Handle(c *gin.Context) (any, error) {
	_, err := models.CreateUser(service.Username, service.Password)
	return nil, err
}

type GetUserCoursesService struct {
	ID uint `uri:"id" binding:"required"`
}

func (service *GetUserCoursesService) Handle(c *gin.Context) (any, error) {
	user, err := models.GetUserByID(service.ID)
	if err != nil {
		return nil, err
	}
	return user.GetCourses()
}

type GetUserNotifications struct {
	ID uint `uri:"id" binding:"required"`
}

type Notifications struct {
	TeachingHomeworkListsToFinish  []models.Homework `json:"homeworkInProgress"`
	TeachingHomeworkListsToComment []models.Homework `json:"commentInProgress"`
	LeaningHomeworkListsToFinish   []models.Homework `json:"homeworksToBeCompleted"`
	LeaningHomeworkListsToComment  []models.Homework `json:"commentToBeCompleted"`
}

// 返回应该尚未提交的作业,待批阅的作业和每门课最新发布的作业
func (service *GetUserNotifications) Handle(c *gin.Context) (any, error) {
	user, err := models.GetUserByID(service.ID)
	if err != nil {
		return nil, err
	}
	courses, err := user.GetCourses()
	if err != nil {
		return nil, err
	}
	var notifications Notifications
	//得到教的课中进行中和批阅中的作业
	println("len of homework%d", len(courses.LearningCourses))
	//得到学的课中还没完成的作业和还没批阅的作业
	for _, course := range courses.LearningCourses {
		homeworks, err := course.GetHomeworkLists()
		if homeworks == nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(homeworks); j++ {
			if homeworks[j].CommentEndDate.After(time.Now()) {
				if homeworks[j].BeginDate.Before(time.Now()) {
					if homeworks[j].EndDate.After(time.Now()) {
						homework := models.GetHomeWorkSubmissionByHomeworkIDAndUserID(homeworks[j].ID, user.ID)
						if homework == nil {
							notifications.LeaningHomeworkListsToFinish =
								append(notifications.LeaningHomeworkListsToFinish, homeworks[j])
						}
					} else {
						comments, err := models.GetCommentListsByUserIDAndHomeworkID(user.ID, homeworks[j].ID)
						if err != nil {
							return nil, err
						}
						for i := 0; i < len(comments); i++ {
							if comments[i].Grade == -1 {
								notifications.LeaningHomeworkListsToComment =
									append(notifications.TeachingHomeworkListsToComment, homeworks[j])
								break
							}
						}
					}
				}
			}
		}
		for _, course := range courses.TeachingCourses {
			homeworks, err := course.GetHomeworkLists()
			if homeworks == nil {
				continue
			}
			if err != nil {
				return nil, err
			}
			for j := 0; j < len(homeworks); j++ {
				if homeworks[j].CommentEndDate.After(time.Now()) {
					if homeworks[j].BeginDate.Before(time.Now()) {
						if homeworks[j].EndDate.After(time.Now()) {
							notifications.TeachingHomeworkListsToFinish =
								append(notifications.TeachingHomeworkListsToFinish, homeworks[j])
						} else {
							notifications.TeachingHomeworkListsToComment =
								append(notifications.TeachingHomeworkListsToComment, homeworks[j])
						}
					}
				}
			}
		}
	}
	return notifications, nil
}
