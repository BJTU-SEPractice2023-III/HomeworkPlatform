package service

import (
	"errors"
	"fmt"
	"homework_platform/internal/models"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type HomeworkDetail struct {
	CourseID   int `form:"courseid"`
	HomeworkID int `form:"homeworkid"`
}

func (service *HomeworkDetail) Handle(c *gin.Context) (any, error) {
	homework, err2 := models.GetHomeworkByID(uint(service.HomeworkID))
	if err2 != nil {
		return nil, errors.New("没有找到该作业")
	}
	path := fmt.Sprintf("./data/homeworkassign/%d/%d", service.CourseID, service.HomeworkID)
	files, err := os.ReadDir(path)
	homework.FilePaths = make([]string, 0)
	if err == nil {
		for _, file := range files {
			filePath := filepath.Join(path, file.Name())
			homework.FilePaths = append(homework.FilePaths, filePath)
		}
	}
	return homework, nil
}

type AssignHomeworkService struct {
	CourseID    int                     `form:"courseid"`
	Name        string                  `form:"name"`
	Description string                  `form:"description"`
	BeginDate   time.Time               `form:"begindate"`
	EndDate     time.Time               `form:"enddate"`
	Files       []*multipart.FileHeader `form:"files"`
}

func (service *AssignHomeworkService) Handle(c *gin.Context) (any, error) {
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

type HomeworkLists struct {
	CourseID int `form:"courseid"`
}

func (service *HomeworkLists) Handle(c *gin.Context) (any, error) {
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

type DeleteHomework struct {
	CourseID   int `form:"courseid"`
	HomeworkID int `form:"homeworkid"`
}

func (service *DeleteHomework) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能删除不是您的课程的作业")
	}
	homework, err2 := models.GetHomeworkByID(uint(service.HomeworkID))
	if err2 != nil {
		return nil, err2
	}
	if homework.CourseID != service.CourseID {
		return nil, errors.New("该作业并非属于该课程")
	}
	if err := homework.Deleteself(); err != nil {
		return nil, err
	}
	dirPath := fmt.Sprintf("./data/homeworkassign/%d/%d", service.CourseID, service.HomeworkID)
	os.RemoveAll(dirPath)

	return nil, nil
}

type UpdateHomeworkService struct {
	CourseID    int                     `form:"courseid"`
	HomeworkID  int                     `form:"homeworkid"`
	Name        string                  `form:"name"`
	Description string                  `form:"description"`
	BeginDate   time.Time               `form:"begindate"`
	EndDate     time.Time               `form:"enddate"`
	Files       []*multipart.FileHeader `form:"files"`
}

func (service *UpdateHomeworkService) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能查看不是您的课程的作业")
	}
	//CourseID
	homework, err2 := models.GetHomeworkByID(uint(service.HomeworkID))
	if err2 != nil {
		return nil, errors.New("没有找到该作业")
	}
	if homework.CourseID != service.CourseID {
		return nil, errors.New("不能更改不是对应课程的作业")
	}
	homework.UpdateInformation(service.Name, service.Description, service.BeginDate, service.EndDate)
	os.RemoveAll(fmt.Sprintf("./data/homeworkassign/%d/%d", service.CourseID, homework.ID))

	for _, f := range service.Files {
		log.Println(f.Filename)
		dst := fmt.Sprintf("./data/homeworkassign/%d/%d/%s", service.CourseID, homework.ID, f.Filename)
		// 上传文件到指定的目录
		c.SaveUploadedFile(f, dst)
	}
	println(service.CourseID)
	return nil, nil
}

type SubmitListsService struct {
	CourseID   int `form:"courseid"`
	HomeworkID int `form:"homeworkid"`
}

func (service *SubmitListsService) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能查看不是您的课程的作业")
	}
	//CourseID
	homework, err2 := models.GetHomeworkByIDWithSubmissionLists(uint(service.HomeworkID))
	if err2 != nil {
		return nil, errors.New("没有找到该作业")
	}
	return homework.HomeworkSubmissions, nil
}
