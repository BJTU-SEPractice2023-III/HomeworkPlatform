package models

import (
	"errors"
	"fmt"
	"homework_platform/internal/utils"
	"log"
	"math"
	"strings"
	"time"
)

type User struct {
	// gorm.Model
	ID        uint   `json:"id" gorm:"primaryKey"`
	Username  string `json:"username" gorm:"unique; not null"`                                                                       // 用户名
	Password  string `json:"-" gorm:"not null"`                                                                                      // 密码
	IsAdmin   bool   `json:"isAdmin"`                                                                                                // 是否是管理员
	Signature string `json:"signature"`                                                                                              // 用户个性签名
	Avatar    string `json:"avatar" gorm:"default:'https://s1.imagehub.cc/images/2023/12/12/36100f1b1b03d8170712fc8a4dc49e4b.jpeg'"` // 用户头像
	//// Associations ////
	// A user has many courses
	// Also check course.go
	// Check: https://gorm.io/docs/has_many.html
	TeachingCourses []*Course `json:"teachingCourses" gorm:"foreignKey:TeacherID;constraint:OnDelete:CASCADE"` //引用了Course这个字段作为外键

	// A student has many courses, a course has many students
	// Also check course.go
	// Check: https://gorm.io/docs/many_to_many.html
	LearningCourses []*Course `json:"learningCourses" gorm:"many2many:user_courses;constraint:OnDelete:CASCADE"`

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

// CreateUser creates a user with the given username and raw password
// Tested
func CreateUser(username string, password string) (*User, error) {
	logPrefix := fmt.Sprintf("[models/user]: CreateUser(username: %s)", username)

	log.Printf("%s: 正在创建...", logPrefix)
	password = utils.EncodePassword(password, utils.RandStringRunes(16))
	user := User{Username: username, Password: password, IsAdmin: false} // 默认创建的用户权限为普通用户

	if err := DB.Create(&user).Error; err != nil {
		log.Printf("%s: 创建失败(%s)", logPrefix, err)
		return nil, err
	}
	log.Printf("%s: 创建成功(id = %v)", logPrefix, user.ID)
	return &user, nil
}

// GetUserById gets the user corresponding to the given id
// Tested
func GetUserByID(id uint) (user User, err error) {
	logPrefix := fmt.Sprintf("[models/user]: GetUserByID(id: %d)", id)

	log.Printf("%s: 正在查找...", logPrefix)
	if err = DB.Preload("LearningCourses").Preload("TeachingCourses").First(&user, id).Error; err != nil {
		log.Printf("%s: 查找失败(%s)", logPrefix, err)
	} else {
		log.Printf("%s: 查找成功(username = %s)", logPrefix, user.Username)
	}
	return
}

// DeleteUserById deletes the user corresponding to the given id
// Tested
func DeleteUserById(id uint) (err error) {
	logPrefix := fmt.Sprintf("[models/user]: DeleteUserByID(id: %d)", id)

	log.Printf("%s: 正在删除...", logPrefix)
	if err = DB.Delete(&User{}, id).Error; err != nil {
		log.Printf("%s: 删除失败(%s)", logPrefix, err)
	} else {
		log.Printf("%s: 删除成功", logPrefix)
	}
	return
}

// GetUsers gets all users
// Tested
func GetUsers() (users []User, err error) {
	logPrefix := "[models/user]: GetUsers"

	log.Printf("%s: 正在获取...", logPrefix)
	if err = DB.Find(&users).Error; err != nil {
		log.Printf("%s: 获取失败", logPrefix)
	} else {
		log.Printf("%s: 获取完成(len = %d)", logPrefix, len(users))
	}

	return users, nil
}

// GetFiles get files the user owned
// Tested in file.go
func (user *User) GetFiles() ([]File, error) {
	logPrefix := fmt.Sprintf("[models/user]: (*User<id: %d>).GetFiles", user.ID)

	log.Printf("%s: 正在获取文件...", logPrefix)
	var files []File
	err := DB.Model(user).Association("Files").Find(&files)
	if err != nil {
		log.Printf("%s: 获取失败(%s)", logPrefix, err)
		return nil, err
	}
	log.Printf("%s: 获取成功(len = %d)", logPrefix, len(files))
	return files, nil
}

func (user *User) ChangeAvatar(url string) error {
	if url == "" {
		return errors.New("图床传入错误")
	}
	if err := DB.Model(&user).Updates(User{Avatar: url}).Error; err != nil {
		return err
	}
	return nil
}

// CheckPassword checks whether the password is correct or not
func (user *User) CheckPassword(password string) bool {
	salt := strings.Split(user.Password, ":")[0]
	return user.Password == utils.EncodePassword(password, salt)
}

// ChangePassword changes password
func (user *User) ChangePassword(password string) error {
	if len(password) == 0 {
		return errors.New("密码不能为空")
	}
	if err := DB.Model(&user).Updates(User{Password: utils.EncodePassword(password, utils.RandStringRunes(16))}).Error; err != nil {
		return err
	}
	return nil
}

// CreateCourse creates a course
// Tested in course_test.go
func (user *User) CreateCourse(name string, begindate time.Time, enddate time.Time, description string) (*Course, error) {
	logPrefix := fmt.Sprintf("[models/user]: (*User<id: %d>).CreateCourse(nam: %s)", user.ID, name)

	log.Printf("%s: 正在创建...", logPrefix)
	course := Course{
		Name:        name,
		BeginDate:   begindate,
		EndDate:     enddate,
		Description: description,
		TeacherID:   user.ID,
	}
	if err := DB.Create(&course).Error; err != nil {
		return nil, err
	}

	return &course, nil
}

// SelectCourse selects a course
// Tested in course_test.go
func (user *User) SelectCourse(courseId uint) error {
	course, err := GetCourseByID(courseId)
	if err != nil {
		return err
	}
	res := course.GetStudentsByID(user.ID)
	if res {
		return errors.New("无法重复选课")
	}
	if err = DB.Model(user).Association("LearningCourses").Append(&course); err != nil {
		return err
	}
	return nil
}

//// Old things ////

func (user *User) UpdateDegree(averageGrade int, myGrade int) error {
	// 如何更新,上限为700,下限为300
	diff := math.Abs(float64(averageGrade - myGrade))
	// 更新动量为当abs()
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
