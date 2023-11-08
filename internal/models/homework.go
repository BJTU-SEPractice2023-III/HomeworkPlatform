package models

import (
	"errors"
	"fmt"
	"homework_platform/internal/bootstrap"
	"io/ioutil"
	"log"
	"path/filepath"
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
	// Also check homeworkSubmission.go
	// Check: https://gorm.io/docs/has_many.html
	HomeworkSubmissions []HomeworkSubmission `json:"-"`
	FilePaths           []string             `json:"file_paths" gorm:"-"`
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

func (homeworkd Homework) Deleteself() error {
	res := DB.Delete(&homeworkd)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func CreateHomework(id uint, name string, description string,
	begindate time.Time, endtime time.Time, commentendate time.Time) (any, error) {
	log.Printf("正在创建作业<id=%d>", id)
	if begindate.After(endtime) {
		return nil, errors.New("结束时间不可早于开始时间")
	}
	if endtime.After(commentendate) {
		return nil, errors.New("评论开始时间不可早于结束时间")
	}
	log.Printf("homework<name:%s>:正在创建作业", name)
	if name == "" {
		log.Printf("homework<name:%s>:作业名字不可为空", name)
		return nil, errors.New("名称不可为空")
	}
	if description == "" {
		log.Printf("homework<name:%s>:作业内容不可为空", name)
		return nil, errors.New("内容不可为空")
	}
	if bootstrap.Sqlite {
		//TODO:这里好像sqlite不会生成外键约束
		_, err := GetCourseByID(id)
		if err != nil {
			return nil, errors.New("课程不存在")
		}
	}
	newhomework := Homework{
		CourseID:       id,
		Name:           name,
		Description:    description,
		BeginDate:      begindate,
		EndDate:        endtime,
		CommentEndDate: commentendate,
	}
	res := DB.Create(&newhomework)
	if res.Error != nil {
		log.Printf("homework<name:%s>:创建作业失败", name)
		return nil, errors.New("创建失败")
	}
	log.Printf("homework<name:%s>:创建作业成功,作业id为%d", name, newhomework.ID)
	return newhomework, nil
}

func GetHomeworkByID(id uint) (Homework, error) {
	log.Printf("正在查找<Homework>(ID = %d)...", id)
	var work Homework
	res := DB.First(&work, id)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return work, res.Error
	}
	root := fmt.Sprintf("./data/homeworkassign/%d/", work.ID)
	files, err := ioutil.ReadDir(root)
	if err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			work.FilePaths = append(work.FilePaths, filepath.Join(root, file.Name()))
		}
	}
	log.Printf("查找完成: <Homeworkd>(homeworkName = %s)", work.Name)
	return work, nil
}

func GetHomeworkByIDWithSubmissionLists(id uint) (Homework, error) {
	log.Printf("正在查找<Homework>(ID = %d)...", id)
	var work Homework

	res := DB.Preload("HomeworkSubmissions").First(&work, id)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return work, res.Error
	}
	log.Printf("查找完成: <Homeworkd>(homeworkName = %s)", work.Name)
	return work, nil
}
