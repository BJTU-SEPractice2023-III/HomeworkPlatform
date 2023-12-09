package service

import (
	"errors"
	"fmt"
	"homework_platform/internal/models"
	"log"
	"mime/multipart"

	// "net/http"
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
	id := c.GetUint("ID")
	user, err := models.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if service.BeginDate.After(service.EndDate) {
		return nil, errors.New("开始时间晚于结束时间")
	}
	course, err := user.CreateCourse(service.Name, service.BeginDate, service.EndDate, service.Description)
	return course.ID, err
}

type UpdateCourseDescription struct {
	CourseID    uint   `uri:"id" binding:"required"`
	Description string `form:"description"`
}

func (service *UpdateCourseDescription) Handle(c *gin.Context) (any, error) {
	var err error
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
	CourseID uint `uri:"id" binding:"required"`
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
	// id, _ := c.Get("ID")
	// if course.TeacherID != id {
	// 	return nil, errors.New("不能查看不是您的课程的学生列表")
	// }
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
	userId := c.GetUint("ID")
	user, err := models.GetUserByID(userId)
	if err != nil {
		return nil, err
	}
	if err = user.SelectCourse(service.CourseID); err != nil {
		return nil, err
	}
	return nil, nil
}

// type GetTeachingCourse struct {
// }

// func (service *GetTeachingCourse) Handle(c *gin.Context) (any, error) {
// 	id, _ := c.Get("ID")
// 	user, err := models.GetUserByID(id.(uint))
// 	if err != nil {
// 		return nil, err
// 	}
// 	courses, res := user.GetTeachingCourse()
// 	if res != nil {
// 		return nil, res
// 	}
// 	return courses, nil
// }

// type GetLearningCourse struct{}

// func (service *GetLearningCourse) Handle(c *gin.Context) (any, error) {
// 	id, _ := c.Get("ID")
// 	user, err := models.GetUserByID(id.(uint))
// 	if err != nil {
// 		return nil, err
// 	}
// 	courses, res := user.GetLearningCourse()
// 	if res != nil {
// 		return nil, res
// 	}
// 	return courses, nil
// }

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

type StudentHomework struct {
	models.Homework
	Submitted bool `json:"submitted"`
	Score     int  `json:"score"`
}

type GetCourseHomeworks struct {
	CourseID uint `uri:"id" binding:"required"`
}

func (service *GetCourseHomeworks) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}
	homeworks, err := course.GetHomeworks()
	if err != nil {
		return nil, err
	}

	id := c.GetUint("ID")
	if id != course.TeacherID {
		studentHomeworks := []StudentHomework{}
		for _, homework := range homeworks {
			studentHomework := StudentHomework{
				Homework:  homework,
				Submitted: false,
				Score:     -1,
			}

			homeworkSubmission, err := homework.GetSubmissionByUserId(id)
			if homeworkSubmission != nil || err != nil{
				studentHomework.Submitted = true
				studentHomework.Score = homeworkSubmission.Score
			}

			studentHomeworks = append(studentHomeworks, studentHomework)
		}
		return studentHomeworks, nil
	} else {
		// TODO: Frontend oriented 的老师需要的作业信息（比如额外增加提交人数统计等）
		return homeworks, nil
	}
}

// This is a special one using no bindings (manualy do the bindings in Handle func)
type CreateCourseHomework struct {
	CourseID       uint                    `uri:"id" binding:"required"`
	Name           string                  `form:"name"`
	Description    string                  `form:"description"`
	BeginDate      time.Time               `form:"beginDate"`
	EndDate        time.Time               `form:"endDate"`
	CommentEndDate time.Time               `form:"commentEndDate"`
	Files          []*multipart.FileHeader `form:"files"`
}

func (s *CreateCourseHomework) Handle(c *gin.Context) (any, error) {
	if c.ContentType() != "multipart/form-data" {
		return nil, errors.New("not supported content-type")
	}

	var err error
	// 从 Uri 获取 CourseID
	err = c.ShouldBindUri(s)
	if err != nil {
		return nil, err
	}
	// 从 Form 获取其他数据
	err = c.ShouldBind(s)
	if err != nil {
		return nil, err
	}
	log.Print(s)
	if s.Name == "" {
		return nil, errors.New("名称不能为空")
	}
	if len(s.Files) == 0 && s.Description == "" {
		log.Printf("作业没有内容")
		return nil, errors.New("内容不能为空")
	}
	if s.BeginDate.After(s.EndDate) || s.EndDate.After(s.CommentEndDate) {
		log.Printf("时间混乱")
		return nil, errors.New("时间顺序错误")
	}
	// 获取课程
	course, err := models.GetCourseByID(s.CourseID)
	if err != nil {
		return nil, err
	}
	// 获取请求者 ID
	id, _ := c.Get("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能发布不是您的课程的作业")
	}
	// 创建课程
	homework, err := course.CreateHomework(
		s.Name,
		s.Description,
		s.BeginDate,
		s.EndDate,
		s.CommentEndDate,
	)
	if err != nil {
		return nil, errors.New("创建失败")
	}
	for _, f := range s.Files {
		file, err := models.CreateFileFromFileHeaderAndContext(f, c)
		if err != nil {
			// TODO: err handle
		} else {
			file.Attach(homework.ID, models.TargetTypeHomework)
		}
	}
	return homework.ID, nil
}
