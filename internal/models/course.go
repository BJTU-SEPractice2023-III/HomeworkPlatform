package models

import (
	"errors"
	"time"
)

type Course struct {
	// gorm.Model
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	BeginDate   time.Time `json:"beginDate"`
	EndDate     time.Time `json:"endDate"`
	Description string    `json:"description" gorm:"not null"`

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

	Complaints []Complaint `josn:"-" gorm:"constraint:OnDelete:CASCADE"`
}

// CreateHomework creates a homework
// Tested in homework_test.go
func (course Course) CreateHomework(name string, description string, begindate time.Time, endtime time.Time, commentendate time.Time) (*Homework, error) {
	// // logPrefix := fmt.Sprintf("[models/course]: (Course<id: %d>).CreateHomework<name: %s>", course.ID, name)

	// log.Printf("%s: 正在创建...", // logPrefix)
	if begindate.After(endtime) {
		// log.Printf("%s: 结束时间不可早于开始时间\n", // logPrefix)
		return nil, errors.New("结束时间不可早于开始时间")
	}
	if endtime.After(commentendate) {
		// log.Printf("%s: 评论结束时间不可早于作业结束时间\n", // logPrefix)
		return nil, errors.New("评论结束时间不可早于作业结束时间")
	}

	homework := Homework{
		CourseID:       course.ID,
		Name:           name,
		Description:    description,
		BeginDate:      begindate,
		EndDate:        endtime,
		CommentEndDate: commentendate,
	}
	if err := DB.Create(&homework).Error; err != nil {
		// log.Printf("%s: 创建失败(%s)\n", // logPrefix, err)
		return nil, err
	}
	// log.Printf("%s: 创建成功(id = %d)\n", // logPrefix, homework.ID)
	return &homework, nil
}

// GetCourseByID gets course by id
// Tested
func GetCourseByID(id uint) (Course, error) {
	// // logPrefix := fmt.Sprintf("[models/course]: GetCourseByID(id: %d)", id)
	// log.Printf("%s: 正在查找...", // logPrefix)
	var course Course

	res := DB.First(&course, id)
	if res.Error != nil {
		// log.Printf("%s: 查找失败: %s", // logPrefix, res.Error)
		return course, res.Error
	}
	// log.Printf("%s: 查找完成(CourseName = %s)", // logPrefix, course.Name)
	return course, nil
}

// GetHomeworks gets all homeworks of a course
// Tested
func (course Course) GetHomeworks() ([]Homework, error) {
	var homeworks []Homework
	err := DB.Model(&course).Preload("Files").Preload("HomeworkSubmissions").Association("Homeworks").Find(&homeworks)
	if err != nil {
		return nil, err
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
	var students []*User
	err := DB.Model(&course).Association("Students").Find(&students, "id = ?", id)
	if err != nil {
		// log.Println(err)
		return false
	}

	// log.Println(students)
	return len(students) > 0
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

// func CreateCourse(name string, begindate time.Time,
// 	enddate time.Time, description string, teachderID uint) (uint, error) {
// 	// log.Printf("正在创建<Course>(name = %s)", name)
// 	if len(name) == 0 {
// 		// log.Printf("创建失败,课程名不能为空")
// 		return 0, errors.New("创建失败")
// 	}
// 	if enddate.Before(begindate) {
// 		// log.Printf("创建失败,开始时间不能晚于结束时间")
// 		return 0, errors.New("创建失败")
// 	}
// 	c := Course{
// 		Name:        name,
// 		BeginDate:   begindate,
// 		EndDate:     enddate,
// 		Description: description,
// 		TeacherID:   uint(teachderID),
// 	}
// 	res := DB.Create(&c)
// 	if res.Error != nil {
// 		return 0, errors.New("创建失败")
// 	}

// 	return c.ID, nil
// }

func (course Course) Deleteself() error {
	res := DB.Delete(&course)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
