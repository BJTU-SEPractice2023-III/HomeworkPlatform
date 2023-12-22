package models

import (
	"errors"
	"fmt"
	"homework_platform/internal/utils"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type File struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	// Size in Bytes
	Size           uint   `json:"size"`
	Path           string `json:"path"`
	UserID         uint   `json:"userID"`
	AttachmentID   uint   `json:"AttachmentID"`
	AttachmentType string `json:"AttachmentType"`
}

// Tested
func createFile(userId uint, name string, size uint, path string) (*File, error) {
	// logPrefix := fmt.Sprintf("[models/file]: CreateFile<name: %s, size: %d, path: %s>", name, size, path)

	// log.Printf("%s: 正在创建...", // logPrefix)
	file := File{
		UserID: userId,
		Name:   name,
		Size:   size,
		Path:   path,
	}
	res := DB.Create(&file)
	if res.Error != nil {
		// log.Printf("%s: 创建失败(%s)", // logPrefix, res.Error)
		return nil, res.Error
	}
	// log.Printf("%s: 创建成功(id = %d)", // logPrefix, file.ID)
	return &file, nil
}

// CreateFileFromFileHeaderAndContext save the file to filesystem,
// and create a file record based on c.GetUint("id")
// ! Not tested yet, test it with router
func CreateFileFromFileHeaderAndContext(fileHeader *multipart.FileHeader, c *gin.Context) (*File, error) {
	// logPrefix := "[models/file]: CreateFileFromFileHeaderAndContext"

	id := c.GetUint("ID")
	file := File{
		UserID: id,
		Name:   fileHeader.Filename,
		Size:   uint(fileHeader.Size),
		Path:   fmt.Sprintf("./data/%d/%s-%s", id, utils.GetTimeStamp(), fileHeader.Filename),
	}
	// log.Printf("%s: 正在保存文件到文件系统(path: %s)...", // logPrefix, file.Path)
	if err := c.SaveUploadedFile(fileHeader, file.Path); err != nil {
		// log.Printf("%s: 保存失败(%s)", // logPrefix, err)
		return nil, err
	}
	return createFile(id, file.Name, file.Size, file.Path)
}

// Tested
func GetFileByID(id uint) (File, error) {
	// logPrefix := fmt.Sprintf("[models/file]: GetFileById<id: %d>", id)

	// log.Printf("%s: 正在查找...", // logPrefix)
	var file File
	res := DB.First(&file, id)
	if res.Error != nil {
		// log.Printf("%s: 查找失败: %s", // logPrefix, res.Error)
		return file, res.Error
	}
	// log.Printf("%s: 查找成功: <File>(name = %s)", // logPrefix, file.Name)
	return file, nil
}

// Tested
func DeleteFileById(id uint) error {
	// logPrefix := fmt.Sprintf("[models/file]: DeleteFileById<id: %d>", id)

	// log.Printf("%s: 正在删除...", // logPrefix)
	res := DB.Delete(&File{}, id)
	if res.Error != nil {
		// log.Printf("%s: 删除失败(%s)\n", // logPrefix, res.Error)
		return res.Error
	}
	// log.Printf("%s: 删除成功(id = %d)\n", // logPrefix, id)
	return nil
}

const (
	TargetTypeHomework = iota
	TargetTypeHomeworkSubmission
)

// Attach to homework or homeworkSubmission
func (file *File) Attach(id uint, targetType uint) error {
	if targetType == TargetTypeHomework {
		homework, err := GetHomeworkByID(id)
		if err != nil {
			return err
		}
		_, err = homework.addAttachment(file)
		if err != nil {
			return err
		}
	} else if targetType == TargetTypeHomeworkSubmission {
		homeworkSubmission, err := GetHomeworkSubmissionById(id)
		if err != nil {
			return err
		}
		_, err = homeworkSubmission.addAttachment(file)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Unknown attach type")
	}
	return nil
}
