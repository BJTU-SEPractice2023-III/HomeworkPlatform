package models

import (
	"errors"
	"homework_platform/internal/utils"
	"log"
	"math"
	"strings"
	// "gorm.io/gorm"
)

type User struct {
	// gorm.Model
	ID        uint   `json:"id" gorm:"primaryKey"`
	Username  string `json:"username" gorm:"unique"` // 用户名
	Password  string `json:"-"`                      // 密码
	IsAdmin   bool   `json:"isAdmin"`                // 是否是管理员
	Signature string `json:"signature"`              // 用户个性签名
	////// Associations //////
	// A user has many courses
	// Also check course.go
	// Check: https://gorm.io/docs/has_many.html
	TeachingCourses []*Course `json:"-" gorm:"foreignKey:TeacherID;constraint:OnDelete:CASCADE"` //引用了Course这个字段作为外键

	// A student has many courses, a course has many students
	// Also check course.go
	// Check: https://gorm.io/docs/many_to_many.html
	LearningCourses []*Course `json:"-" gorm:"many2many:user_courses;constraint:OnDelete:CASCADE"`

	// A user has many homework submissions
	// Also check homework_submissions.go
	// Check: https://gorm.io/docs/has_many.html
	HomeworkSubmissions []HomeworkSubmission `json:"-" gorm:"constraint:OnDelete:CASCADE"`

	// A user has many comments
	// Also check comment.go
	// Check: https://gorm.io/docs/has_many.html
	Comments []Comment `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Files    []File    `json:"-" gorm:"constraint:OnDelete:CASCADE"`

	Complaints []Complaint `josn:"-" gorm:"constraint:OnDelete:CASCADE"`
	// 算法设计,根据置信度的比率来打分和,在根据均值的偏差计算置信度
	DegreeOfConfidence float64 `json:"-" gorm:"default:300"`
}

func (user *User) UploadFile(name string, size uint, path string) (uint, error) {
	file := File {
		Name: name,
		Size: size,
		Path: path,
	}
	err := DB.Model(user).Association("Files").Append(&file)
	if err != nil {
		return 0, err
	}
	return file.ID, nil
}

func (user *User) GetFiles() ([]File, error) {
	var files []File
	err := DB.Model(user).Association("Files").Find(&files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (user *User) CheckPassword(password string) bool {
	salt := strings.Split(user.Password, ":")[0]
	log.Printf("用户密码为: %s", user.Password)
	log.Printf("盐: %s", salt)

	return user.Password == utils.EncodePassword(password, salt)
}

func (user *User) UpdateDegree(averageGrade int, myGrade int) error {
	//如何更新,上限为700,下限为300
	diff := math.Abs(float64(averageGrade - myGrade))
	//更新动量为当abs()
	if diff < 5 {
		user.DegreeOfConfidence += -2*diff + 10
	} else {
		user.DegreeOfConfidence -= 10*diff/95 - 50.0/95.0
	}
	if user.DegreeOfConfidence > 700 {
		user.DegreeOfConfidence = 700
	} else if user.DegreeOfConfidence < 300 {
		user.DegreeOfConfidence = 300
	}
	result := DB.Model(&user).Updates(User{DegreeOfConfidence: user.DegreeOfConfidence})
	return result.Error
}

func UpgradeToAdmin(userId uint) error {
	user, err := GetUserByID(userId)
	if err != nil {
		return err
	}
	user.IsAdmin = true
	err = DB.Save(&user).Error
	return err
}

func (user *User) ChangeSignature(signature string) error {
	log.Printf("正在修改签名<User>(Username = %s, Signature = %s)...", user.Username, signature)
	result := DB.Model(&user).Updates(User{Signature: signature})
	return result.Error
}

func (user *User) ChangePassword(password string) bool {
	log.Printf("正在修改密码<User>(Username = %s, Password = %s)...", user.Username, password)
	if len(password) == 0 {
		log.Printf("修改失败,用户密码不能为空")
		return false
	}
	result := DB.Model(&user).Updates(User{Password: utils.EncodePassword(password, utils.RandStringRunes(16))})
	return result.Error == nil
}

func (user *User) GetTeachingCourse() ([]*Course, error) {
	res := DB.Preload("TeachingCourses").First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user.TeachingCourses, nil
}

func (user *User) GetLearningCourse() ([]*Course, error) {
	res := DB.Preload("LearningCourses").First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user.LearningCourses, nil
}

func (user *User) DeleteSelf() bool {
	log.Printf("正在删除用户<User>(Username = %s)...", user.Username)
	res := DB.Delete(&user)
	return res.Error == nil
}

func CreateUser(username string, password string) (uint, error) {
	log.Printf("正在创建<User>(Username = %s, Password = %s)...", username, password)
	if len(username) == 0 {
		return 0, errors.New("名称不能为空")
	}
	if len(password) == 0 {
		return 0, errors.New("密码不能为空")
	}
	password = utils.EncodePassword(password, utils.RandStringRunes(16))
	user := User{Username: username, Password: password, IsAdmin: false} //默认创建的用户权限为普通用户

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

const (
	Learning = iota
	Teaching
)

type UserCourse struct {
	TeachingCourses []*Course `json:"teachingCourses"`
	LearningCourses []*Course `json:"learningCourses"`
}

func (user *User) GetCourses() (UserCourse, error) {
	var res UserCourse
	var err error

	if res.TeachingCourses, err = user.GetTeachingCourse(); err != nil {
		return res, nil
	}
	if res.LearningCourses, err = user.GetLearningCourse(); err != nil {
		return res, nil
	}
	return res, nil
}
