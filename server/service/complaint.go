package service

import (
	"errors"
	"homework_platform/internal/models"
	"log"

	"github.com/gin-gonic/gin"
)

type CreateComplaint struct {
	Reason     string `form:"reason"`
	HomeworkID uint   `uri:"id" binding:"required"`
}

func (service *CreateComplaint) Handle(c *gin.Context) (any, error) {
	//绑定id
	err := c.ShouldBindUri(service)
	if err != nil {
		return nil, err
	}
	//绑定reason
	err = c.ShouldBind(service)
	if err != nil {
		return nil, err
	}

	homework, err := models.GetHomeworkByID(service.HomeworkID)
	if err != nil {
		return nil, err
	}

	log.Printf("正在创建Complaint<homeworkId:%d,reason:%s>", service.HomeworkID, service.Reason)
	id := c.GetUint("ID")
	if service.Reason == "" {
		return nil, errors.New("原因不能为空")
	}
	homework_submission, err := homework.GetSubmissionByUserId(id)
	if err != nil {
		return nil, errors.New("没有找到该提交")
	}
	log.Printf("作业的用户为id%d", homework_submission.UserID)
	if homework_submission.UserID != id {
		return nil, errors.New("这不是您的作业")
	}

	complaint, _ := models.GetComplaintByHomeworkIdAndUserId(homework.ID, id)
	if complaint != nil && complaint.ID != 0 {
		return nil, errors.New("未查询")
	}
	err = models.CreateTeacherComplaint(
		homework_submission.ID,
		homework_submission.HomeworkID,
		homework.CourseID,
		service.Reason,
	)
	return nil, err
}

type DeleteComplaint struct {
	ComplaintId uint `uri:"id" binding:"required"`
}

func (service *DeleteComplaint) Handle(c *gin.Context) (any, error) {
	err := models.DeleteComplaint(service.ComplaintId)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type UpdateComplaint struct {
	ComplaintId uint   `uri:"id" binding:"required"`
	Reason      string `form:"reason"`
}

func (service *UpdateComplaint) Handle(c *gin.Context) (any, error) {
	//绑定id
	err := c.ShouldBindUri(service)
	if err != nil {
		return nil, err
	}
	//绑定reason
	err = c.ShouldBind(service)
	if err != nil {
		return nil, err
	}
	if service.Reason == "" {
		return nil, errors.New("原因不能为空")
	}
	complain, err := models.GetComplaintById(service.ComplaintId)
	if err != nil {
		return nil, err
	}
	log.Printf("修改reason为:%s", service.Reason)
	complain.Reason = service.Reason
	err = complain.Save()
	return nil, err
}

type SolveComplaint struct {
	ComplaintId uint `uri:"id" binding:"required"`
}

func (service *SolveComplaint) Handle(c *gin.Context) (any, error) {
	complaint, err := models.GetComplaintById(service.ComplaintId)
	if err != nil {
		return nil, err
	}
	complaint.Solved = true
	err = complaint.Save()
	return nil, err
}

type GetComplaint struct {
	HomeworkID uint `uri:"id" binding:"required"`
}

func (service *GetComplaint) Handle(c *gin.Context) (any, error) {
	id := c.GetUint("ID")
	homework, err := models.GetHomeworkByID(service.HomeworkID)
	if err != nil {
		return nil, err
	}
	course, err := models.GetCourseByID(homework.CourseID)
	if err != nil {
		return nil, err
	}
	if course.TeacherID == id {
		complaints, err := models.GetComplaintByHomeworkID(service.HomeworkID)
		if err != nil {
			return nil, err
		}
		return complaints, nil
	} else {
		submission, err := homework.GetSubmissionByUserId(id)
		if submission == nil || err != nil {
			return nil, errors.New("未提交作业")
		}
		complaint, err := models.GetComplaintBySubmissionID(submission.ID)
		if err != nil {
			return nil, err
		}
		return complaint, nil
	}
}
