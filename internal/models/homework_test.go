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

	user, _ := CreateUser("username", "password")
	course, _ := user.CreateCourse("course", time.Now(), time.Now().AddDate(1, 0, 0), "desc")

	date := time.Date(2023, 12, 1, 22, 40, 0, 0, time.UTC)
	var homeworkData Homework
	var homework *Homework
	var err error

	// INFO: 成功
	homeworkData = Homework{
		Name:           "homework",
		Description:    "desc",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, 7),
		CommentEndDate: date.AddDate(0, 0, 14),
	}

	homework, _ = course.CreateHomework(
		homeworkData.Name,
		homeworkData.Description,
		homeworkData.BeginDate,
		homeworkData.EndDate,
		homeworkData.CommentEndDate,
	)

	res, err := GetHomeworkByID(homework.ID)
	assert.Nil(err)
	assert.Equal(homework.ID, res.ID)
	assert.Equal(course.ID, res.CourseID)
	assert.Equal(homeworkData.Name, res.Name)
	assert.Equal(homeworkData.Description, res.Description)
	assert.Equal(homeworkData.BeginDate, res.BeginDate)
	assert.Equal(homeworkData.EndDate, res.EndDate)
	assert.Equal(homeworkData.CommentEndDate, res.CommentEndDate)

	// INFO: 失败：结束时间早于开始时间
	homeworkData = Homework{
		Name:           "homework",
		Description:    "desc",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, -7),
		CommentEndDate: date.AddDate(0, 0, 14),
	}
	_, err = course.CreateHomework(
		homeworkData.Name,
		homeworkData.Description,
		homeworkData.BeginDate,
		homeworkData.EndDate,
		homeworkData.CommentEndDate,
	)
	assert.Error(err)

	// INFO: 失败：评论结束时间早于作业结束时间
	homeworkData = Homework{
		Name:           "homework",
		Description:    "desc",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, 77),
		CommentEndDate: date,
	}
	_, err = course.CreateHomework(
		homeworkData.Name,
		homeworkData.Description,
		homeworkData.BeginDate,
		homeworkData.EndDate,
		homeworkData.CommentEndDate,
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
	_, err = course.CreateHomework(
		homeworkData.Name,
		homeworkData.Description,
		homeworkData.BeginDate,
		homeworkData.EndDate,
		homeworkData.CommentEndDate,
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
	_, err = course.CreateHomework(
		homeworkData.Name,
		homeworkData.Description,
		homeworkData.BeginDate,
		homeworkData.EndDate,
		homeworkData.CommentEndDate,
	)
	assert.Error(err)
}

func TestGetHomeworkById(t *testing.T) {
	deleteData()
	assert := assert.New(t)

	user, _ := CreateUser("username", "password")
	course, _ := user.CreateCourse("course", time.Now(), time.Now().AddDate(1, 0, 0), "desc")

	date := time.Date(2023, 12, 1, 22, 40, 0, 0, time.UTC)
	var homeworkData Homework
	var err error

	// INFO: 成功
	homeworkData = Homework{
		CourseID:       course.ID,
		Name:           "homework",
		Description:    "desc",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, 7),
		CommentEndDate: date.AddDate(0, 0, 14),
	}

	homework, _ := course.CreateHomework(
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

	user, _ := CreateUser("username", "password")
	course, _ := user.CreateCourse("course", time.Now(), time.Now().AddDate(1, 0, 0), "desc")

	date := time.Date(2023, 12, 1, 22, 40, 0, 0, time.UTC)
	var homeworkData Homework
	var err error

	// INFO: 成功
	homeworkData = Homework{
		CourseID:       course.ID,
		Name:           "homework",
		Description:    "desc",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, 7),
		CommentEndDate: date.AddDate(0, 0, 14),
	}

	homework, _ := course.CreateHomework(
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

	user, _ := CreateUser("username", "password")
	course, _ := user.CreateCourse("course", time.Now(), time.Now().AddDate(1, 0, 0), "desc")

	date := time.Date(2023, 12, 1, 22, 40, 0, 0, time.UTC)
	var homeworkData Homework
	var err error

	// INFO: 成功
	homeworkData = Homework{
		CourseID:       course.ID,
		Name:           "homework",
		Description:    "desc",
		BeginDate:      date,
		EndDate:        date.AddDate(0, 0, 7),
		CommentEndDate: date.AddDate(0, 0, 14),
	}

	homework, _ := course.CreateHomework(
		homeworkData.Name,
		homeworkData.Description,
		homeworkData.BeginDate,
		homeworkData.EndDate,
		homeworkData.CommentEndDate,
	)

	file := File{
		UserID: user.ID,
		Name:   "name3",
		Size:   3,
		Path:   "./data/2/testfile3",
	}

	attachment, err := homework.addAttachment(&file)
	assert.Nil(err)
	assert.Equal(file.UserID, attachment.UserID)
	assert.Equal(file.Name, attachment.Name)
	assert.Equal(file.Size, attachment.Size)
	assert.Equal(file.Path, attachment.Path)


	homework, err = GetHomeworkByID(homework.ID)
	assert.Nil(err)
	assert.Equal(homework.Files[0], *attachment)
}

