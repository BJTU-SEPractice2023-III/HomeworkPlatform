package models

import (
	"log"

	"gorm.io/gorm"
)

type TeacherNotice struct {
	gorm.Model

	UserID               uint   `form:"userId"`
	TeacherID            uint   `form:"teacherId"`
	HomeworkSubmissionID uint   `form:"homeworkSubmissionId"`
	HomeworkID           uint   `form:"homeworkId"`
	CourseID             uint   `form:"courseId"`
	Solved               bool   `form:"solved"`
	Reason               string `form:"reason"`
}

func CreateTeacherNotice(submissionId uint, homeworkId uint, CourseID uint, reason string) error {
	log.Printf("正在创建<Notice>(SubmissionId = %d)...", submissionId)
	homeworkSubmission := GetHomeWorkSubmissionByID(submissionId)
	course, err := GetCourseByID(CourseID)
	if err != nil {
		return err
	}
	notice := TeacherNotice{HomeworkSubmissionID: submissionId,
		HomeworkID: homeworkId,
		CourseID:   CourseID,
		Reason:     reason,
		UserID:     homeworkSubmission.ID,
		TeacherID:  course.TeacherID,
	}
	notice.Solved = false
	res := DB.Create(&notice)
	if res.Error == nil {
		log.Printf("创建完成<Notice>(ID = %v)", notice.ID)
	}
	return res.Error
}

func DeleteNotice(noticeId uint) error {
	log.Printf("正在删除修改请求<Notice>(id = %d)...", noticeId)
	notice, err := GetNoticeById(noticeId)
	if err != nil {
		return err
	}
	res := DB.Delete(&notice)
	return res.Error
}

func SolveNotice(noticeID uint) error {
	notice, err := GetNoticeById(noticeID)
	if err != nil {
		return err
	}
	notice.Solved = true
	res := DB.Save(&notice)
	return res.Error
}

func GetNoticeById(Id uint) (TeacherNotice, error) {
	log.Printf("正在查找<Notice>(ID = %d)...", Id)
	var notice TeacherNotice

	res := DB.Where("homework_submission_id=?", Id).First(&notice)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return notice, res.Error
	}
	return notice, nil
}

func GetNoticeBySubmission(submissionId uint) (TeacherNotice, error) {
	log.Printf("正在查找<Notice>(submissionId = %d)...", submissionId)
	var notice TeacherNotice

	res := DB.Where("homework_submission_id=?", submissionId).First(&notice)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return notice, res.Error
	}
	return notice, nil
}

func GetNoticeByUserID(UserID uint) ([]TeacherNotice, error) {
	log.Printf("正在查找<Notice>(userId = %d)...", UserID)
	var notice []TeacherNotice

	res := DB.Where("user_id=?", UserID).Find(&notice)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return notice, res.Error
	}
	return notice, nil
}

func GetNoticeByTeacherID(teacherId uint) ([]TeacherNotice, error) {
	log.Printf("正在查找<Notice>(TeacherID = %d)...", teacherId)
	var notice []TeacherNotice

	res := DB.Where("teacher_id=?", teacherId).Find(&notice)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return notice, res.Error
	}
	return notice, nil
}
