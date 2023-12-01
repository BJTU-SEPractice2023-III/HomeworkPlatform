package models

import (
	"testing"
	"time"

	// "time"
	"github.com/stretchr/testify/assert"
)

func TestCreateHomework(t *testing.T) {
	deleteData()
	assert := assert.New(t)

	userId, _ := CreateUser("username", "password")
	courseId, _ := CreateCourse("course", time.Now(), time.Now().AddDate(1, 0, 0), "desc", userId)

	date := time.Date(2023, 12, 1, 22, 40, 0, 0, time.UTC)
	var homeworkData Homework
	var homework *Homework
	var err error

	// INFO: 成功
	homeworkData = Homework{
		CourseID:       courseId,
		Name:           "homework",
		Description:    "desc",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, 7),
		CommentEndDate: date.AddDate(0, 0, 14),
	}

	homework, _ = CreateHomework(
		homeworkData.CourseID,
		homeworkData.Name,
		homeworkData.Description,
		homeworkData.BeginDate,
		homeworkData.EndDate,
		homeworkData.CommentEndDate,
	)

	res, err := GetHomeworkByID(homework.ID)
	assert.Nil(err)
	assert.Equal(homework.ID, res.ID)
	assert.Equal(homeworkData.CourseID, res.CourseID)
	assert.Equal(homeworkData.Name, res.Name)
	assert.Equal(homeworkData.Description, res.Description)
	assert.Equal(homeworkData.BeginDate, res.BeginDate)
	assert.Equal(homeworkData.EndDate, res.EndDate)
	assert.Equal(homeworkData.CommentEndDate, res.CommentEndDate)

	// INFO: 失败：不存在的 CourseID
	homeworkData = Homework{
		CourseID:       0,
		Name:           "homework",
		Description:    "desc",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, 7),
		CommentEndDate: date.AddDate(0, 0, 14),
	}
	_, err = CreateHomework(
		homework.CourseID,
		homework.Name,
		homework.Description,
		homework.BeginDate,
		homework.EndDate,
		homework.CommentEndDate,
	)
	assert.Error(err)

	// INFO: 失败：结束时间早于开始时间
	homeworkData = Homework{
		CourseID:       0,
		Name:           "homework",
		Description:    "desc",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, -7),
		CommentEndDate: date.AddDate(0, 0, 14),
	}
	_, err = CreateHomework(
		homework.CourseID,
		homework.Name,
		homework.Description,
		homework.BeginDate,
		homework.EndDate,
		homework.CommentEndDate,
	)
	assert.Error(err)

	// INFO: 失败：评论结束时间早于作业结束时间
	homeworkData = Homework{
		CourseID:       0,
		Name:           "homework",
		Description:    "desc",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, 77),
		CommentEndDate: date,
	}
	_, err = CreateHomework(
		homework.CourseID,
		homework.Name,
		homework.Description,
		homework.BeginDate,
		homework.EndDate,
		homework.CommentEndDate,
	)
	assert.Error(err)

	// INFO: 失败：名称为空
	homeworkData = Homework{
		CourseID:       0,
		Name:           "",
		Description:    "desc",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, 77),
		CommentEndDate: date,
	}
	_, err = CreateHomework(
		homework.CourseID,
		homework.Name,
		homework.Description,
		homework.BeginDate,
		homework.EndDate,
		homework.CommentEndDate,
	)
	assert.Error(err)

	// INFO: 失败：描述为空
	homeworkData = Homework{
		CourseID:       0,
		Name:           "name",
		Description:    "",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, 77),
		CommentEndDate: date,
	}
	_, err = CreateHomework(
		homework.CourseID,
		homework.Name,
		homework.Description,
		homework.BeginDate,
		homework.EndDate,
		homework.CommentEndDate,
	)
	assert.Error(err)
}

func TestGetHomeworkById(t *testing.T) {
	deleteData()
	assert := assert.New(t)

	userId, _ := CreateUser("username", "password")
	courseId, _ := CreateCourse("course", time.Now(), time.Now().AddDate(1, 0, 0), "desc", userId)

	date := time.Date(2023, 12, 1, 22, 40, 0, 0, time.UTC)
	var homeworkData Homework
	var err error

	// INFO: 成功
	homeworkData = Homework{
		CourseID:       courseId,
		Name:           "homework",
		Description:    "desc",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, 7),
		CommentEndDate: date.AddDate(0, 0, 14),
	}

	homework, _ := CreateHomework(
		homeworkData.CourseID,
		homeworkData.Name,
		homeworkData.Description,
		homeworkData.BeginDate,
		homeworkData.EndDate,
		homeworkData.CommentEndDate,
	)

	res, err := GetHomeworkByID(homework.ID)
	assert.Nil(err)
	assert.Equal(homework.ID, res.ID)
	assert.Equal(homeworkData.CourseID, res.CourseID)
	assert.Equal(homeworkData.Name, res.Name)
	assert.Equal(homeworkData.Description, res.Description)
	assert.Equal(homeworkData.BeginDate, res.BeginDate)
	assert.Equal(homeworkData.EndDate, res.EndDate)
	assert.Equal(homeworkData.CommentEndDate, res.CommentEndDate)
}

func TestDeleteHomeworkById(t *testing.T) {
	deleteData()
	assert := assert.New(t)

	userId, _ := CreateUser("username", "password")
	courseId, _ := CreateCourse("course", time.Now(), time.Now().AddDate(1, 0, 0), "desc", userId)

	date := time.Date(2023, 12, 1, 22, 40, 0, 0, time.UTC)
	var homeworkData Homework
	var err error

	// INFO: 成功
	homeworkData = Homework{
		CourseID:       courseId,
		Name:           "homework",
		Description:    "desc",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, 7),
		CommentEndDate: date.AddDate(0, 0, 14),
	}

	homework, _ := CreateHomework(
		homeworkData.CourseID,
		homeworkData.Name,
		homeworkData.Description,
		homeworkData.BeginDate,
		homeworkData.EndDate,
		homeworkData.CommentEndDate,
	)

	err = DeleteHomeworkById(homework.ID)
	assert.Nil(err)

	_, err = GetHomeworkByID(homework.ID)
	assert.Error(err)
}

func TestHomeworkAttachment(t *testing.T) {
	deleteData()
	assert := assert.New(t)

	userId, _ := CreateUser("username", "password")
	courseId, _ := CreateCourse("course", time.Now(), time.Now().AddDate(1, 0, 0), "desc", userId)

	date := time.Date(2023, 12, 1, 22, 40, 0, 0, time.UTC)
	var homeworkData Homework
	var err error

	// INFO: 成功
	homeworkData = Homework{
		CourseID:       courseId,
		Name:           "homework",
		Description:    "desc",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, 7),
		CommentEndDate: date.AddDate(0, 0, 14),
	}

	homework, _ := CreateHomework(
		homeworkData.CourseID,
		homeworkData.Name,
		homeworkData.Description,
		homeworkData.BeginDate,
		homeworkData.EndDate,
		homeworkData.CommentEndDate,
	)

	file := File{
		UserID: userId,
		Name:   "name3",
		Size:   3,
		Path:   "./data/2/testfile3",
	}

	attachment, err := homework.AddAttachment(file.UserID, file.Name, file.Size, file.Path)
	assert.Nil(err)
	assert.Equal(file.UserID, attachment.UserID)
	assert.Equal(file.Name, attachment.Name)
	assert.Equal(file.Size, attachment.Size)
	assert.Equal(file.Path, attachment.Path)


	homework, err = GetHomeworkByID(homework.ID)
	assert.Nil(err)
	assert.Equal(homework.Files[0], *attachment)
}
