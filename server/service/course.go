package service

import (
	"errors"
	"fmt"
	"homework_platform/internal/models"
	"log"
	"mime/multipart"
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
	if service.BeginDate.After(service.EndDate) {
		return nil, errors.New("开始时间晚于结束时间")
	}
	id, err := models.CreateCourse(service.Name, service.BeginDate, service.EndDate, service.Description, id.(uint))
	return id, err
}

type UpdateCourseDescription struct {
	CourseID    uint   `form:"courseid"`
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
	CourseID uint `form:"courseid"`
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

type GetCourseStudents struct {
	CourseID uint `uri:"id" binding:"required"`
}

func (service *GetCourseStudents) Handle(c *gin.Context) (any, error) {
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

type AddCourseStudentService struct {
	CourseID uint `uri:"id" binding:"required"`
}

func (service *AddCourseStudentService) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}

	id, _ := c.Get("ID")
	// 查看该用户是否已经选择了course
	res := course.GetStudentsByID(id.(uint))
	if res {
		return nil, errors.New("无法重复选课")
	}

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

type GetCourseStudentLists struct {
	CourseID uint `form:"courseid"`
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
	CourseID uint `form:"courseid"`
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

type GetCourses struct{}

func (service *GetCourses) Handle(c *gin.Context) (any, error) {
	courses, err := models.GetCourses()
	return courses, err
}

type GetCourse struct {
	ID uint `uri:"id" binding:"required"`
}

func (service *GetCourse) Handle(c *gin.Context) (any, error) {
	fmt.Println(*service)
	course, err := models.GetCourseByID(service.ID)
	return course, err
}

type GetCourseHomeworks struct {
	CourseID uint `uri:"id" binding:"required"`
}

func (service *GetCourseHomeworks) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能查看不是您的课程的作业")
	}
	homeworks, err2 := course.GetHomeworkLists()
	if err2 != nil {
		return nil, err2
	}
	return homeworks, nil
}

// TODO: 有问题，需要把 uri 和 form 分割为两个结构体进行两次 bind，需要重构下框架。
type AddCourseHomework struct {
	CourseID       uint                    `uri:"id" binding:"required"`
	Name           string                  `form:"name"`
	Description    string                  `form:"description"`
	BeginDate      time.Time               `form:"begindate"`
	EndDate        time.Time               `form:"enddate"`
	CommentEndDate time.Time               `form:"commentenddate"`
	Files          []*multipart.FileHeader `form:"files"`
}

func (service *AddCourseHomework) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能发布不是您的课程的作业")
	}
	//CourseID
	homework, err2 := models.CreateHomework(
		service.CourseID,
		service.Name,
		service.Description,
		service.BeginDate,
		service.EndDate,
		service.CommentEndDate,
	)
	if err2 != nil {
		return nil, errors.New("创建失败")
	}
	for _, f := range service.Files {
		log.Println(f.Filename)
		dst := fmt.Sprintf("./data/homeworkassign/%d/%d/%s", service.CourseID, homework.(models.Homework).ID, f.Filename)
		// 上传文件到指定的目录
		c.SaveUploadedFile(f, dst)
	}
	return nil, nil
}