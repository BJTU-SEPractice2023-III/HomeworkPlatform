package models

import (
	"errors"
	"fmt"
	"log"
	"time"
	"gorm.io/gorm"
)

type Homework struct {
	gorm.Model
	CourseID       uint      `json:"courseId" gorm:"type:int(20)"`
	Name           string    `json:"name" gorm:"type:varchar(255)"`
	Description    string    `json:"description"`
	BeginDate      time.Time `json:"beginDate"`
	EndDate        time.Time `json:"endDate"`
	CommentEndDate time.Time `json:"commentEndDate"`
	Assigned       int       `json:"-" gorm:"default:-1"`
	// A homework has many submissions
	// Also check homework_submission.go
	// Check: https://gorm.io/docs/has_many.html
	HomeworkSubmissions []HomeworkSubmission `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	FilePaths           []string             `json:"file_paths" gorm:"-"`
	Files               []File               `json:"-" gorm:"constraint:OnDelete:CASCADE; polymorphic:Attachment;"`
}

// Tested
func CreateHomework(courseId uint, name string, description string,
	begindate time.Time, endtime time.Time, commentendate time.Time) (*Homework, error) {
	logPrefix := fmt.Sprintf("[models/homework]: CreateHomework<courseId: %d, name: %s>", courseId, name)

	log.Printf("%s: 正在创建...", logPrefix)
	if begindate.After(endtime) {
		log.Printf("%s: 结束时间不可早于开始时间\n", logPrefix)
		return nil, errors.New("结束时间不可早于开始时间")
	}
	if endtime.After(commentendate) {
		log.Printf("%s: 评论结束时间不可早于作业结束时间\n", logPrefix)
		return nil, errors.New("评论结束时间不可早于作业结束时间")
	}
	if name == "" {
		log.Printf("%s: 名称不可为空\n", logPrefix)
		return nil, errors.New("名称不可为空")
	}
	if description == "" {
		log.Printf("%s: 内容不可为空\n", logPrefix)
		return nil, errors.New("内容不可为空")
	}

	homework := Homework{
		CourseID:       courseId,
		Name:           name,
		Description:    description,
		BeginDate:      begindate,
		EndDate:        endtime,
		CommentEndDate: commentendate,
	}
	res := DB.Create(&homework)
	if res.Error != nil {
		log.Printf("%s: 创建失败(%s)\n", logPrefix, res.Error)
		return nil, res.Error
	}
	log.Printf("%s: 创建成功(id = %d)\n", logPrefix, homework.ID)
	return &homework, nil
}

// Tested
func DeleteHomeworkById(id uint) error {
	logPrefix := fmt.Sprintf("[models/homework]: DeleteHomeworkById<id: %d>", id)

	log.Printf("%s: 正在删除...", logPrefix)
	res := DB.Delete(&Homework{}, id)
	if res.Error != nil {
		log.Printf("%s: 删除失败(%s)\n", logPrefix, res.Error)
		return res.Error
	}
	log.Printf("%s: 删除成功(id = %d)\n", logPrefix, id)
	return nil
}

// Tested
func GetHomeworkByID(id uint) (*Homework, error) {
	logPrefix := fmt.Sprintf("[models/homework]: GetHomeworkById<id: %d>", id)

	log.Printf("%s: 正在查找...", logPrefix)
	var homework Homework
	res := DB.Preload("HomeworkSubmissions").Preload("Files").First(&homework, id)
	if res.Error != nil {
		log.Printf("%s: 查找失败: %s", logPrefix, res.Error)
		return &homework, res.Error
	}
	// homework.GetFiles()
	log.Printf("%s: 查找成功: <Homework>(name = %s)", logPrefix, homework.Name)
	return &homework, nil
}

// Tested
func (homework *Homework) AddAttachment(userId uint, name string, size uint, path string) (*File, error) {
	file := File {
		UserID: userId,
		Name: name,
		Size: size,
		Path: path,
	}
	err := DB.Model(homework).Association("Files").Append(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// TODO: split the argument into each fields
func (homework *Homework) AddSubmission(submission HomeworkSubmission) (uint, error) {
	logPrefix := fmt.Sprintf("[models/homework]: homework<id: %d>.AddSubmission<userId: %d>", homework.ID, submission.UserID)

	log.Printf("%s: 正在创建...", logPrefix)
	res := DB.Create(&submission)
	if res.Error != nil {
		log.Printf("%s: 创建失败(%s)", logPrefix, res.Error)
		return 0, res.Error
	}
	log.Printf("%s: 创建成功(id = %d)", logPrefix, submission.ID)
	return submission.ID, nil
}

func (homework *Homework) UpdateInformation(name string, desciption string, beginDate time.Time, endDate time.Time, commentendate time.Time) bool {
	log.Printf("正在修改homework<id:%d>的详细信息", homework.ID)
	if beginDate.After(endDate) {
		log.Printf("homework<id:%d>:开始时间不可晚于结束时间", homework.ID)
		return false
	}
	if endDate.After(commentendate) {
		log.Printf("homework<id:%d>:结束时间不可晚于批阅时间", homework.ID)
		return false
	}
	if name == "" {
		log.Printf("homework<id:%d>:作业名字不可为空", homework.ID)
		return false
	}
	if desciption == "" {
		log.Printf("homework<name:%d>:作业内容不可为空", homework.ID)
		return false
	}
	result := DB.Model(&homework).Updates(Homework{
		Name: name, Description: desciption, BeginDate: beginDate, EndDate: endDate, CommentEndDate: commentendate,
	})
	return result.Error == nil
}

func GetSubmittedUsers(id uint) ([]User, error) {
	homework, err := GetHomeworkByID(id)
	if err != nil {
		return nil, err
	}

	users := []User{}
	for _, submissions := range homework.HomeworkSubmissions {
		user, err := GetUserByID(submissions.UserID)
		if err == nil {
			users = append(users, user)
		}
	}

	return users, nil
}
