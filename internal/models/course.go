package models

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Name        string    `json:"name"`
	BeginDate   time.Time `json:"begin_date"`
	EndDate     time.Time `json:"end_date"`
	Description string    `json:"description"`

	// A teacher has many course
	// Also check user.go
	// Check: https://gorm.io/docs/has_many.html
	TeacherID uint

	// A student has many Course, a course has many students
	// Also check user.go
	// Check: https://gorm.io/docs/many_to_many.html
	Students []*User `gorm:"many2many:user_courses;"`

	// A course has many homework
	// Also check homework.go
	// Check: https://gorm.io/docs/has_many.html
	Homeworks []Homework
}

func GetAllStudents(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Model(&User{}).Preload("Course").Find(&users).Error
	return users, err
}

func (course Course) UpdateCourseDescription(description string) bool {
	result := DB.Model(&course).Updates(Course{Description: description})
	if result.Error != nil {
		return false
	}
	return true
}

func CreateCourse(name string, begindate time.Time,
	enddate time.Time, description string, teachderID uint) error {
	c := Course{
		Name:        name,
		BeginDate:   begindate,
		EndDate:     enddate,
		Description: description,
		TeacherID:   uint(teachderID),
	}
	res := DB.Create(&c)
	if res.Error != nil {
		return errors.New("创建失败")
	}

	return nil
}

func GetCourseByID(id int) (Course, error) {
	log.Printf("正在查找<Course>(ID = %d)...", id)
	var course Course

	res := DB.First(&course, id)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return course, res.Error
	}
	log.Printf("查找完成: <Course>(CourseName = %s)", course.Name)
	return course, nil
}

func (course Course) Deleteself() error {
	res := DB.Delete(&course)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
