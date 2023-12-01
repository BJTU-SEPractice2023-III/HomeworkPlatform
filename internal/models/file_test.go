package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateFile(t *testing.T) {
	assert := assert.New(t)

	var file *File
	var files []File
	var err error

	userId, _ := CreateUser("teacher", "password")
	user, _ := GetUserByID(userId)
	
	// INFO: No Attach
	fileData := File {
		UserID: userId,
		Name: "name",
		Size: 1,
		Path: "./data/1/testfile",
	}
	// fileId, err := createFile(file.UserID, file.Name, file.Size, file.Path)
	fileId, err := user.UploadFile(fileData.Name, fileData.Size, fileData.Path)
	assert.Nil(err)

	res, err := GetFileByID(fileId)
	assert.Nil(err)
	assert.Equal(fileId, res.ID)
	assert.Equal(fileData.Name, res.Name)
	assert.Equal(fileData.Size, res.Size)
	assert.Equal(fileData.Path, res.Path)

	files, err = user.GetFiles()
	assert.Nil(err)
	fmt.Println(files)

	courseId, _ := CreateCourse("course", time.Now(), time.Now().AddDate(1, 0, 0), "desc", userId)
	// course, _ := GetCourseByID(courseId)

	homework, _ := CreateHomework(courseId, "homework", "desc", time.Now(), time.Now().AddDate(0, 0, 7), time.Now().AddDate(0, 0, 14))

	// INFO: Attach type - homeworks
	fileData = File {
		UserID: userId,
		Name: "name2",
		Size: 2,
		Path: "./data/1/testfile2",
	}
	file, err = homework.AddAttachment(fileData.UserID, fileData.Name, fileData.Size, fileData.Path)
	assert.Nil(err)
	res, err = GetFileByID(file.ID)
	assert.Nil(err)
	assert.Equal(file.ID, res.ID)
	assert.Equal(fileData.UserID, res.UserID)
	assert.Equal(fileData.Name, res.Name)
	assert.Equal(fileData.Size, res.Size)
	assert.Equal(fileData.Path, res.Path)

	files, err = user.GetFiles()
	assert.Nil(err)
	for _, file := range files {
		fmt.Println(file)
	}

	// INFO: Attach tupe - homework_submissions
	studentId, _ := CreateUser("student", "password")
	student, _ := GetUserByID(studentId)	

	SelectCourse(studentId, courseId)

	homeworkSubmissionData := HomeworkSubmission {
		HomeworkID: homework.ID,
		UserID: studentId,
		Content: "content",
	}
	submissionId, _ := homework.AddSubmission(homeworkSubmissionData)
	homeworkSubmission := GetHomeWorkSubmissionByID(submissionId)

	fileData = File {
		UserID: studentId,
		Name: "name3",
		Size: 3,
		Path: "./data/2/testfile3",
	}

	fileId, err = homeworkSubmission.AddAttachment(fileData.UserID, fileData.Name, fileData.Size, fileData.Path)
	assert.Nil(err)
	res, err = GetFileByID(fileId)
	assert.Nil(err)
	assert.Equal(fileId, res.ID)
	assert.Equal(fileData.UserID, res.UserID)
	assert.Equal(fileData.Name, res.Name)
	assert.Equal(fileData.Size, res.Size)
	assert.Equal(fileData.Path, res.Path)

	files, err = homeworkSubmission.GetAttachments()
	assert.Nil(err)
	assert.Equal(len(files), 1)
	assert.Equal(files[0].ID, uint(3))

	files, err = student.GetFiles()
	assert.Nil(err)
	assert.Equal(len(files), 1)
	assert.Equal(files[0].ID, uint(3))

	for _, file := range files {
		fmt.Println(file)
	}
}

func TestDeleteFileById(t *testing.T) {
	assert := assert.New(t)

	err := DeleteFileById(3)
	assert.Nil(err)

	user, _ := GetUserByID(2)
	files, err := user.GetFiles()
	assert.Nil(err)
	for file := range files {
		fmt.Println(file)
	}
}