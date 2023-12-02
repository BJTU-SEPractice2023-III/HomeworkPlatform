package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCourseGetHomeworks(t *testing.T) {
	deleteData()
	assert := assert.New(t)

	user, _ := CreateUser("username", "password")
	course, _ := user.CreateCourse("course", time.Now(), time.Now().AddDate(0, 0, 7), "desc")

	homework1, _ := course.CreateHomework("homework1", "desc", time.Now(), time.Now().AddDate(0, 0, 7), time.Now().AddDate(0, 0, 14))
	homework1, _ = GetHomeworkByID(homework1.ID)
	homework2, _ := course.CreateHomework("homework2", "desc", time.Now(), time.Now().AddDate(0, 0, 7), time.Now().AddDate(0, 0, 14))
	homework2, _ = GetHomeworkByID(homework2.ID)
	homework3, _ := course.CreateHomework("homework3", "desc", time.Now(), time.Now().AddDate(0, 0, 7), time.Now().AddDate(0, 0, 14))
	homework3, _ = GetHomeworkByID(homework3.ID)

	homeworks, err := course.GetHomeworks()
	assert.Nil(err)
	assert.Equal([]Homework{*homework1, *homework2, *homework3}, homeworks)
}
