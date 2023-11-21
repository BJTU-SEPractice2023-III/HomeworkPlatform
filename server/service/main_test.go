package service_test

import (
	"encoding/json"
	"homework_platform/internal/bootstrap"
	"homework_platform/internal/config"
	"homework_platform/internal/models"
	"homework_platform/server"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Router *gin.Engine

type ResponseData struct {
	Token string `json:"token"`
}

// 解码 JSON 响应
var responseData struct {
	Data ResponseData `json:"data"`
}

func GetAuthorziation(w *httptest.ResponseRecorder) string {
	responseBody := w.Body.Bytes()

	err := json.Unmarshal(responseBody, &responseData)
	if err != nil {
		panic(err)
	}
	// 访问数据
	return responseData.Data.Token
}

func CreateData() {
	models.CreateUser("xyh", "123")
	models.CreateUser("xeh", "123")
	models.CreateUser("xsh", "123")

}

func TestMain(m *testing.M) {
	bootstrap.Sqlite = true
	var err error
	models.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	bootstrap.Sqlite = true
	if err != nil {
		panic(err)
	}
	bootstrap.Config = &config.Config{JWTSigningString: "moorxJ", SQLDSN: "123"}

	models.DB.AutoMigrate(&models.User{})
	models.DB.AutoMigrate(&models.Course{})
	models.DB.AutoMigrate(&models.Homework{})
	models.DB.AutoMigrate(&models.HomeworkSubmission{})
	models.DB.AutoMigrate(&models.Comment{})
	models.DB.AutoMigrate(&models.Complaint{})
	models.DB.Create(&models.User{})
	CreateData()
	Router = server.InitRouter()
	// api.Run(":8888")
	m.Run()
}

// 	api.Run(":8888")
// 	os.Exit(m.Run())
