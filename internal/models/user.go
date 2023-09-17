package models

import (
	"homework_platform/internal/utils"
	"log"
	"strings"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username string `json:"username" gorm:"unique"` // 用户名
	Password string `json:"-"`                      // 密码
	IsAdmin  bool   `json:"is_admin"`               // 是否是管理员

	////// Associations //////
	// A user has many courses
	// Also check course.go
	// Check: https://gorm.io/docs/has_many.html
	TeachingCourses []Course `gorm:"foreignKey:TeacherID"`

	// A student has many courses, a course has many students
	// Also check course.go
	// Check: https://gorm.io/docs/many_to_many.html
	LearningCourses []*Course `gorm:"many2many:user_courses;"`

	// A user has many homework submissions
	// Also check homework_submissions.go
	// Check: https://gorm.io/docs/has_many.html
	HomeworkSubmissions []HomeworkSubmission

	// A user has many comments
	// Also check comment.go
	// Check: https://gorm.io/docs/has_many.html
	Comments []Comment
}

// TODO: implement methods

func (user *User) CheckPassword(password string) bool {
	salt := strings.Split(user.Password, ":")[0]
	log.Printf("用户密码为: %s", user.Password)
	log.Printf("盐: %s", salt)

	return user.Password == utils.EncodePassword(password, salt)
}

func CreateUser(username string, password string) (uint, error) {
	log.Printf("正在创建<User>(Username = %s, Password = %s)...", username, password)
	user := User{Username: username, Password: password}

	res := DB.Create(&user)
	if res.Error == nil {
		log.Printf("创建完成<User>(ID = %v, Username = %s, Password = %s)...", user.ID, user.Username, user.Password)
	}
	return user.ID, res.Error
}


func GetUserByID(id uint) (User, error) {
	log.Printf("正在查找<User>(ID = %d)...", id)
	var user User

	res := DB.First(&user, id)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return user, res.Error
	}
	log.Printf("查找完成: <User>(Username = %s)", user.Username)
	return user, nil
}

func GetUserByUsername(username string) (User, error) {
	log.Printf("正在查找<User>(Username = %s)...", username)
	var user User

	res := DB.Where("username = ?", username).First(&user)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return user, res.Error
	}
	log.Printf("查找完成: <User>(Username = %s)", user.Username)
	return user, nil
}

func GetUserList() ([]User, error) {
	log.Println("正在获取所有 User...")
	var userList = make([]User, 0)

	res := DB.Find(&userList)
	if res.Error != nil {
		log.Printf("获取失败: %s", res.Error)
		return userList, res.Error
	}
	log.Printf("获取完成：共 %d 条数据", len(userList))

	return userList, nil
}
