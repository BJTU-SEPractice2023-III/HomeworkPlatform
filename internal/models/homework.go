package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Homework struct {
	gorm.Model
	CourseID       uint      `json:"courseId"`
	Name           string    `json:"name" gorm:"type:varchar(255)"`
	Description    string    `json:"description"`
	BeginDate      time.Time `json:"beginDate"`
	EndDate        time.Time `json:"endDate"`
	CommentEndDate time.Time `json:"commentEndDate"`
	Assigned       int       `json:"-" gorm:"default:-1"`
	// A homework has many submissions
	// Also check homework_submission.go
	// Check: https://gorm.io/docs/has_many.html
	HomeworkSubmissions []HomeworkSubmission `json:"submissions" gorm:"constraint:OnDelete:CASCADE"`
	// FilePaths           []string             `json:"file_paths" gorm:"-"`
	Files []File `json:"files" gorm:"constraint:OnDelete:CASCADE; polymorphic:Attachment;"`
}

// Tested
func DeleteHomeworkById(id uint) error {
	// logPrefix := fmt.Sprintf("[models/homework]: DeleteHomeworkById<id: %d>", id)

	// log.Printf("%s: 正在删除...", // logPrefix)
	res := DB.Delete(&Homework{}, id)
	if res.Error != nil {
		// log.Printf("%s: 删除失败(%s)\n", // logPrefix, res.Error)
		return res.Error
	}
	// log.Printf("%s: 删除成功(id = %d)\n", // logPrefix, id)
	return nil
}

// Tested
func GetHomeworkByID(id uint) (*Homework, error) {
	// logPrefix := fmt.Sprintf("[models/homework]: GetHomeworkById<id: %d>", id)

	// log.Printf("%s: 正在查找...", // logPrefix)
	var homework Homework
	res := DB.Model(&homework).Preload("HomeworkSubmissions").Preload("Files").First(&homework, id)
	if res.Error != nil {
		// log.Printf("%s: 查找失败: %s", // logPrefix, res.Error)
		return &homework, res.Error
	}
	// homework.GetFiles()
	// log.Printf("%s: 查找成功: <Homework>(name = %s)", // logPrefix, homework.Name)
	return &homework, nil
}

// Tested
func (homework *Homework) addAttachment(file *File) (*File, error) {
	err := DB.Model(homework).Association("Files").Append(file)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// TODO: split the argument into each fields
func (homework *Homework) AddSubmission(userId uint, content string) (*HomeworkSubmission, error) {
	// logPrefix := fmt.Sprintf("[models/homework]: (*Homework<id: %d>).AddSubmission<userId: %d>", homework.ID, userId)

	submission := HomeworkSubmission{
		HomeworkID: homework.ID,
		UserID:     userId,
		Content:    content,
	}

	// log.Printf("%s: 正在创建...", // logPrefix)
	if err := DB.Create(&submission).Error; err != nil {
		// log.Printf("%s: 创建失败(%s)", // logPrefix, err)
		return nil, err
	}
	// log.Printf("%s: 创建成功(id = %d)", // logPrefix, submission.ID)
	return &submission, nil
}

func (homework *Homework) UserCommentFinish(userId uint) (bool, error) {
	_, err := homework.GetSubmissionByUserId(userId)
	if err != nil {
		return false, err
	}
	commens, err := GetCommentsByHomeworkIdAndUserId(homework.ID, userId)
	if err != nil {
		return false, err
	}
	for i := 0; i < len(commens); i++ {
		if commens[i].Score == -1 {
			log.Printf("用户未完成所有批阅")
			return false, nil
		}
	}
	log.Printf("用户完成所有批阅")
	return true, nil
}

// GetSubmissionByUserId gets the submission from a user
func (homework *Homework) GetSubmissionByUserId(userId uint) (*HomeworkSubmission, error) {
	var submission HomeworkSubmission
	if err := DB.Model(&submission).Preload("Files").Where("homework_id = ? AND user_id = ?", homework.ID, userId).First(&submission).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}

// GetSubmissions gets all submissions of a homework
func (homework *Homework) GetSubmissions() ([]HomeworkSubmission, error) {
	var submission []HomeworkSubmission
	if err := DB.Model(homework).Preload("Files").Association("HomeworkSubmissions").Find(&submission); err != nil {
		return nil, err
	}
	return submission, nil
}

// GetSubmissionsWithComments
func (homework *Homework) GetSubmissionsWithComments() ([]HomeworkSubmission, error) {
	var submission []HomeworkSubmission
	if err := DB.Model(homework).Preload("Comments").Association("HomeworkSubmissions").Find(&submission); err != nil {
		return nil, err
	}
	return submission, nil
}

// GetCommentByUserId gets the comment of a user
func (homework *Homework) GetCommentsByUserId(userId uint) ([]Comment, error) {
	var comments []Comment
	if err := DB.Where("homework_id = ? AND user_id = ?", homework.ID, userId).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (homework *Homework) UpdateInformation(name string, desciption string, beginDate time.Time, endDate time.Time, commentendate time.Time) bool {
	// log.Printf("正在修改homework<id:%d>的详细信息", homework.ID)
	if beginDate.After(endDate) {
		// log.Printf("homework<id:%d>:开始时间不可晚于结束时间", homework.ID)
		return false
	}
	if endDate.After(commentendate) {
		// log.Printf("homework<id:%d>:结束时间不可晚于批阅时间", homework.ID)
		return false
	}
	if name == "" {
		// log.Printf("homework<id:%d>:作业名字不可为空", homework.ID)
		return false
	}
	if desciption == "" {
		// log.Printf("homework<name:%d>:作业内容不可为空", homework.ID)
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
