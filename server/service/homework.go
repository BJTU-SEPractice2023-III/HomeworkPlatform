package service

import (
	"errors"
	"fmt"
	"homework_platform/internal/models"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type GetHomeworkById struct {
	ID uint `uri:"id" binding:"required"`
}

func (service *GetHomeworkById) Handle(c *gin.Context) (any, error) {
	homework, err := models.GetHomeworkByID(uint(service.ID))
	if err != nil {
		return nil, err
	}

	course, err := models.GetCourseByID(homework.CourseID)
	if err != nil {
		return nil, err
	}

	id := c.GetUint("ID")
	if id == course.TeacherID {
		return *homework, nil
	} else {
		studentHomework := StudentHomework{
			Homework:  *homework,
			Submitted: false,
			Score:     -1,
		}

		homeworkSubmission, err := homework.GetSubmissionByUserId(id)
		if err == nil {
			studentHomework.Submitted = true
			studentHomework.Score = homeworkSubmission.Score
		}

		return studentHomework, nil
	}
}

type AssignHomeworkService struct {
	CourseID       uint                    `form:"courseId"`
	Name           string                  `form:"name"`
	Description    string                  `form:"description"`
	BeginDate      time.Time               `form:"beginDate"`
	EndDate        time.Time               `form:"endDate"`
	CommentEndDate time.Time               `form:"commentEndDate"`
	Files          []*multipart.FileHeader `form:"files"`
}

func (service *AssignHomeworkService) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}

	id := c.GetUint("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能发布不是您的课程的作业")
	}

	homework, err := course.CreateHomework(
		service.Name,
		service.Description,
		service.BeginDate,
		service.EndDate,
		service.CommentEndDate,
	)
	if err != nil {
		return nil, err
	}
	for _, f := range service.Files {
		file, err := models.CreateFileFromFileHeaderAndContext(f, c)
		if err != nil {
			// TODO: err handle
		} else {
			file.Attach(homework.ID, models.TargetTypeHomework)
		}
	}

	return homework.ID, nil
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
	homeworks, err2 := course.GetHomeworks()
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
	if err := models.DeleteHomeworkById(homework.ID); err != nil {
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
	BeginDate      time.Time               `form:"beginDate"`
	EndDate        time.Time               `form:"endDate"`
	CommentEndDate time.Time               `form:"commentEndDate"`
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

	// TODO: 已经被分配了 需要特殊处理
	if homework.Assigned == 1 {
		//如果提后了提交时间,那么就要重新设置
		if s.EndDate.After(time.Now()) {
			homework.Assigned = -1
			//删除现有评论并且把分数置为0
			err := models.DeleteCommentsByHomeworkID(homework.ID)
			if err != nil {
				return nil, err
			}
			homework_submission := models.GetHomeWorkSubmissionsByHomeworkID(homework.ID)
			for _, submission := range homework_submission {
				submission.Score = -1
				err := submission.UpdateSelf()
				if err != nil {
					return nil, err
				}
			}
			models.DB.Save(&homework)
		}
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

// type SubmitListsService struct {
// 	HomeworkID uint `uri:"id" binding:"required"`
// }

// func (service *SubmitListsService) Handle(c *gin.Context) (any, error) {
// 	query := c.Request.URL.Query()
// 	category := query.Get("all")
// 	if category == "true" {
// 		homework, err2 := models.GetHomeworkByIDWithSubmissionLists(uint(service.HomeworkID))
// 		if err2 != nil {
// 			return nil, errors.New("没有找到该作业")
// 		}
// 		CourseID := homework.CourseID
// 		course, err := models.GetCourseByID(CourseID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		id, _ := c.Get("ID")
// 		if course.TeacherID != id {
// 			return nil, errors.New("不能查看不是您的课程的作业")
// 		}
// 		for i := 0; i < len(homework.HomeworkSubmissions); i++ {
// 			root := fmt.Sprintf("./data/homeworkassign/%d/", homework.HomeworkSubmissions[i].ID)
// 			files, err := os.ReadDir(root)
// 			if err == nil {
// 				for _, file := range files {
// 					if file.IsDir() {
// 						continue
// 					}
// 					homework.HomeworkSubmissions[i].FilePaths = append(homework.HomeworkSubmissions[i].FilePaths, filepath.Join(root, file.Name()))
// 				}
// 			}
// 		}
// 		return homework.HomeworkSubmissions, nil
// 	} else {
// 		id, _ := c.Get("ID")
// 		id = id.(uint)
// 		homework, err := models.GetHomeworkByIDWithSubmissionLists(service.HomeworkID)
// 		if err != nil {
// 			return "该作业号不存在", nil
// 		}
// 		for _, value := range homework.HomeworkSubmissions {
// 			if value.UserID == id {
// 				return value, nil
// 			}
// 		}
// 		return nil, errors.New("该用户未提交作业")
// 	}
// }
