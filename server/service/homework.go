package service

import (
	"errors"
	"fmt"
	"homework_platform/internal/models"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type GetHomework struct {
	ID uint `uri:"id" binding:"required"`
}

func (service *GetHomework) Handle(c *gin.Context) (any, error) {
	homework, err2 := models.GetHomeworkByID(uint(service.ID))
	if err2 != nil {
		return nil, errors.New("没有找到该作业")
	}
	path := fmt.Sprintf("./data/homeworkassign/%d", service.ID)
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
	CourseID       uint                    `form:"courseid"`
	Name           string                  `form:"name"`
	Description    string                  `form:"description"`
	BeginDate      time.Time               `form:"begindate"`
	EndDate        time.Time               `form:"enddate"`
	CommentEndDate time.Time               `form:"commentenddate"`
	Files          []*multipart.FileHeader `form:"files"`
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
		service.CommentEndDate,
	)
	if err2 != nil {
		return nil, errors.New("创建失败")
	}
	for _, f := range service.Files {
		log.Println(f.Filename)
		dst := fmt.Sprintf("./data/homeworkassign/%d/%s", homework.(models.Homework).ID, f.Filename)
		// 上传文件到指定的目录
		c.SaveUploadedFile(f, dst)
	}

	return homework.(models.Homework).ID, nil
}

type HomeworkLists struct {
	CourseID uint `uri:"id" binding:"required"`
}

func (service *HomeworkLists) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}
	// id, _ := c.Get("ID")
	// if course.TeacherID != id {
	// 	return nil, errors.New("不能查看不是您的课程的作业")
	// }
	homeworks, err2 := course.GetHomeworkLists()
	if err2 != nil {
		return nil, err2
	}
	return homeworks, nil
}

type DeleteHomework struct {
	HomeworkID uint `uri:"id" bind:"required"`
}

func (service *DeleteHomework) Handle(c *gin.Context) (any, error) {
	homework, err2 := models.GetHomeworkByID(uint(service.HomeworkID))
	if err2 != nil {
		return nil, err2
	}
	course, err := models.GetCourseByID(homework.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能删除不是您的课程的作业")
	}
	if err := homework.Deleteself(); err != nil {
		return nil, err
	}
	dirPath := fmt.Sprintf("./data/homeworkassign/%d/%d", course.ID, service.HomeworkID)
	os.RemoveAll(dirPath)

	return nil, nil
}

type UpdateHomework struct {
	HomeworkID     uint                    `uri:"id" bind:"required"`
	Name           string                  `form:"name"`
	Description    string                  `form:"description"`
	BeginDate      time.Time               `form:"begindate"`
	EndDate        time.Time               `form:"enddate"`
	CommentEndDate time.Time               `form:"commentenddate"`
	Files          []*multipart.FileHeader `form:"files"`
}

func (s *UpdateHomework) Handle(c *gin.Context) (any, error) {
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
	log.Println(s)

	homework, err := models.GetHomeworkByID(s.HomeworkID)
	if err != nil {
		return nil, err
	}

	course, err := models.GetCourseByID(homework.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能修改不是您的课程的作业")
	}

	homework.UpdateInformation(s.Name, s.Description, s.BeginDate, s.EndDate, s.CommentEndDate)
	os.RemoveAll(fmt.Sprintf("./data/homeworkassign/%d/%d", homework.CourseID, homework.ID))

	for _, f := range s.Files {
		log.Println(f.Filename)
		dst := fmt.Sprintf("./data/homeworkassign/%d/%d/%s", homework.CourseID, homework.ID, f.Filename)
		// 上传文件到指定的目录
		c.SaveUploadedFile(f, dst)
	}
	return nil, nil
}

type SubmitListsService struct {
	HomeworkID uint `uri:"id" binding:"required"`
}

func (service *SubmitListsService) Handle(c *gin.Context) (any, error) {
	query := c.Request.URL.Query()
	category := query.Get("all")
	if category == "true" {
		homework, err2 := models.GetHomeworkByIDWithSubmissionLists(uint(service.HomeworkID))
		if err2 != nil {
			return nil, errors.New("没有找到该作业")
		}
		CourseID := homework.CourseID
		course, err := models.GetCourseByID(CourseID)
		if err != nil {
			return nil, err
		}
		id, _ := c.Get("ID")
		if course.TeacherID != id {
			return nil, errors.New("不能查看不是您的课程的作业")
		}
		for i := 0; i < len(homework.HomeworkSubmissions); i++ {
			root := fmt.Sprintf("./data/homeworkassign/%d/", homework.HomeworkSubmissions[i].ID)
			files, err := ioutil.ReadDir(root)
			if err == nil {
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					homework.HomeworkSubmissions[i].FilePaths = append(homework.HomeworkSubmissions[i].FilePaths, filepath.Join(root, file.Name()))
				}
			}
		}
		return homework.HomeworkSubmissions, nil
	} else {
		id, _ := c.Get("ID")
		id = id.(uint)
		homework, err := models.GetHomeworkByIDWithSubmissionLists(service.HomeworkID)
		if err != nil {
			return "该作业号不存在", nil
		}
		for _, value := range homework.HomeworkSubmissions {
			if value.UserID == id {
				return value, nil
			}
		}
		return nil, errors.New("该用户未提交作业")
	}
}
