package service

import (
	"errors"
	"homework_platform/internal/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateCourse struct {
	Name        string    `form:"name"`
	BeginDate   time.Time `form:"begindate"`
	EndDate     time.Time `form:"enddate"`
	Description string    `form:"description"`
}

func (service *CreateCourse) Handle(c *gin.Context) (any, error) {
	id, _ := c.Get("ID")
	err := models.CreateCourse(service.Name, service.BeginDate, service.EndDate, service.Description, id.(uint))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type UpdateCourseDescription struct {
	CourseID    int    `form:"courseid"`
	Description string `form:"description"`
}

func (service *UpdateCourseDescription) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能修改不是您的课程")
	}
	res := course.UpdateCourseDescription(service.Description)
	if !res {
		return nil, errors.New("创建失败")
	}
	return nil, nil
}

type DeleteCourse struct {
	CourseID int `form:"courseid"`
}

func (service *DeleteCourse) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能修改不是您的课程")
	}
	res := course.Deleteself()
	if res != nil {
		return nil, res
	}
	return nil, nil
}

type GetCourseStudentLists struct {
	CourseID int `form:"courseid"`
}

func (service *GetCourseStudentLists) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能查看不是您的课程的学生列表")
	}
	users, res := course.GetStudents()
	if res != nil {
		return nil, res
	}
	return users, nil
}

type SelectCourseService struct {
	CourseID int `form:"courseid"`
}

func (service *SelectCourseService) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	res := course.GetStudentsByID(id.(uint))
	if res {
		return nil, errors.New("无法重复选课")
	}
	//查看该用户是否已经选择了course
	user, err := models.GetUserByID(id.(uint))
	if err != nil {
		return nil, err
	}
	user.LearningCourses = append(user.LearningCourses, &course)
	result := models.DB.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return nil, nil
}

type GetTeachingCourse struct {
}

func (service *GetTeachingCourse) Handle(c *gin.Context) (any, error) {
	id, _ := c.Get("ID")
	user, err := models.GetUserByID(id.(uint))
	if err != nil {
		return nil, err
	}
	courses, res := user.GetTeachingCourse()
	if res != nil {
		return nil, res
	}
	return courses, nil
}

type GetLearningCourse struct{}

func (service *GetLearningCourse) Handle(c *gin.Context) (any, error) {
	id, _ := c.Get("ID")
	user, err := models.GetUserByID(id.(uint))
	if err != nil {
		return nil, err
	}
	courses, res := user.GetLearningCourse()
	if res != nil {
		return nil, res
	}
	return courses, nil
}

type GetCourses struct {}

func (service *GetCourses) Handle(c *gin.Context) (any, error) {
	courses, err := models.GetCourses()
	log.Println(courses)
	return courses, err
}