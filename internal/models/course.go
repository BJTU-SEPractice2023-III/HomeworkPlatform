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
	BeginDate   time.Time `json:"beginDate"`
	EndDate     time.Time `json:"endDate"`
	Description string    `json:"description"`

	// A teacher has many course
	// Also check user.go
	// Check: https://gorm.io/docs/has_many.html
	TeacherID uint `json:"teacherID"`

	// A student has many Course, a course has many students
	// Also check user.go
	// Check: https://gorm.io/docs/many_to_many.html
	Students []*User `json:"-" gorm:"many2many:user_courses;constraint:OnDelete:CASCADE"`

	// A course has many homework
	// Also check homework.go
	// Check: https://gorm.io/docs/has_many.html
	Homeworks []Homework `json:"-" gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
}

func (course Course) GetHomeworkLists() ([]Homework, error) {
	res := DB.Preload("Homeworks").First(&course, course.ID)
	if res.Error != nil {
		return nil, res.Error
	}
	homeworks := course.Homeworks
	for i := 0; i < len(homeworks); i++ {
		homeworks[i].GetFiles()
	}
	return homeworks, nil
}

func (course Course) GetStudents() ([]*User, error) {
	res := DB.Preload("Students").First(&course)
	if res.Error != nil {
		return nil, res.Error
	}
	return course.Students, nil
}

func (course Course) FindStudents(id uint) bool {
	var student User
	err := DB.Model(&course).Association("Students").Find(&student, "id = ?", id)
	return err == nil
}

func (course Course) SelectCourse(id uint) error {
	res := course.GetStudentsByID(id)
	if res {
		return errors.New("无法重复选课")
	}
	//查看该用户是否已经选择了course
	user, err := GetUserByID(id)
	if err != nil {
		return err
	}
	user.LearningCourses = append(user.LearningCourses, &course)
	result := DB.Save(&user)
	return result.Error

}

func (course Course) UpdateCourseDescription(description string) bool {
	result := DB.Model(&course).Updates(Course{Description: description})
	return result.Error == nil
}

func (course Course) GetStudentsByID(id uint) bool {
	userlists, err := course.GetStudents()
	if err != nil {
		return false
	}
	for _, user := range userlists {
		if user.ID == id {
			return true
		}
	}
	return false
}

func GetCourses() ([]Course, error) {
	var courses []Course
	err := DB.Find(&courses).Error
	return courses, err
}

func CreateCourse(name string, begindate time.Time,
	enddate time.Time, description string, teachderID uint) (uint, error) {
	log.Printf("正在创建<Course>(name = %s)", name)
	if len(name) == 0 {
		log.Printf("创建失败,课程名不能为空")
		return 0, errors.New("创建失败")
	}
	if enddate.Before(begindate) {
		log.Printf("创建失败,开始时间不能晚于结束时间")
		return 0, errors.New("创建失败")
	}
	c := Course{
		Name:        name,
		BeginDate:   begindate,
		EndDate:     enddate,
		Description: description,
		TeacherID:   uint(teachderID),
	}
	res := DB.Create(&c)
	if res.Error != nil {
		return 0, errors.New("创建失败")
	}

	return c.ID, nil
}

func GetCourseByID(id uint) (Course, error) {
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
