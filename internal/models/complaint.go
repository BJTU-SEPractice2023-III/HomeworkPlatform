package models

import (
	"log"

	"gorm.io/gorm"
)

type Complaint struct {
	gorm.Model

	UserID               uint   `form:"userId"`
	TeacherID            uint   `form:"teacherId"`
	HomeworkSubmissionID uint   `form:"homeworkSubmissionId"`
	HomeworkID           uint   `form:"homeworkId"`
	CourseID             uint   `form:"courseId"`
	Solved               bool   `form:"solved"`
	Reason               string `form:"reason"`
}

func GetComplaintByHomeworkIdAndUserId(homeworkId uint, userId uint) (*Complaint, error) {
	log.Printf("正在查找<Complaint>(homeworkId = %d,userId = %d)...", homeworkId, userId)
	var complaint Complaint
	res := DB.Where("homework_id = ? AND user_id = ?", homeworkId, userId).Find(&complaint)
	return &complaint, res.Error
}

func CreateTeacherComplaint(submissionId uint, homeworkId uint, CourseID uint, reason string) error {
	log.Printf("正在创建<Complaint>(SubmissionId = %d)...", submissionId)
	homeworkSubmission := GetHomeWorkSubmissionByID(submissionId)

	course, err := GetCourseByID(CourseID)
	if err != nil {
		return err
	}
	complaint := Complaint{HomeworkSubmissionID: submissionId,
		HomeworkID: homeworkId,
		CourseID:   CourseID,
		Reason:     reason,
		UserID:     homeworkSubmission.ID,
		TeacherID:  course.TeacherID,
	}
	complaint.Solved = false
	res := DB.Create(&complaint)
	if res.Error == nil {
		log.Printf("创建完成<Complaint>(ID = %v)", complaint.ID)
	}
	return res.Error
}

func DeleteComplaint(complaintId uint) error {
	log.Printf("正在删除修改请求<Complaint>(id = %d)...", complaintId)
	complaint, err := GetComplaintById(complaintId)
	if err != nil {
		return err
	}
	res := DB.Delete(&complaint)
	return res.Error
}

func SolveComplaint(complaintID uint) error {
	complaint, err := GetComplaintById(complaintID)
	if err != nil {
		return err
	}
	complaint.Solved = true
	res := DB.Save(&complaint)
	return res.Error
}

func GetComplaintById(Id uint) (Complaint, error) {
	log.Printf("正在查找<Complaint>(ID = %d)...", Id)
	var complaint Complaint

	res := DB.Where("homework_submission_id=?", Id).First(&complaint)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return complaint, res.Error
	}
	return complaint, nil
}

func (complaint *Complaint) Save() error {
	res := DB.Save(&complaint)
	return res.Error
}

func GetComplaintBySubmissionID(submissionId uint) (Complaint, error) {
	log.Printf("正在查找<Complaint>(submissionId = %d)...", submissionId)
	var complaint Complaint

	res := DB.Where("homework_submission_id=?", submissionId).First(&complaint)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return complaint, res.Error
	}
	return complaint, nil
}

func GetComplaintByHomeworkID(homeworkId uint) ([]Complaint, error) {
	log.Printf("正在查找<Complaint>(homeworkId = %d)...", homeworkId)
	var complaints []Complaint

	res := DB.Where("homework_id=?", homeworkId).Find(&complaints)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return complaints, res.Error
	}
	return complaints, nil
}

func GetComplaintByUserID(UserID uint) ([]Complaint, error) {
	log.Printf("正在查找<Complaint>(userId = %d)...", UserID)
	var complaint []Complaint

	res := DB.Where("user_id=? and solved=0", UserID).Find(&complaint)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return complaint, res.Error
	}
	return complaint, nil
}

func GetComplaintByTeacherID(teacherId uint) ([]Complaint, error) {
	log.Printf("正在查找<Complaint>(TeacherID = %d)...", teacherId)
	var complaint []Complaint

	res := DB.Where("teacher_id=? and solved = 0", teacherId).Find(&complaint)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return complaint, res.Error
	}
	return complaint, nil
}
