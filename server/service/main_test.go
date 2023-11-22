package service_test

import (
	"encoding/json"
	"homework_platform/internal/bootstrap"
	"homework_platform/internal/config"
	"homework_platform/internal/models"
	"homework_platform/server"
	"net/http/httptest"
	"testing"
	"time"

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
	models.CreateUser("tjw", "123")
	models.CreateUser("xb", "123")
	models.CreateUser("xbb", "123")
	models.CreateUser("xbbb", "123")

	models.CreateCourse("c++", time.Now(), time.Now().AddDate(0, 0, 1), "c++", 1)
	models.CreateCourse("c3+", time.Now(), time.Now().AddDate(0, 1, 1), "c++", 2)
	models.CreateCourse("c#", time.Now(), time.Now().AddDate(0, 1, 1), "c++", 2)
	models.CreateCourse("c++", time.Now(), time.Now().AddDate(0, 0, 1), "c++", 1)

	models.CreateHomework(1, "c++1", "lkksk", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2))
	models.CreateHomework(1, "c++2", "lkksk", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2))
	models.CreateHomework(2, "c++3", "lkksk", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2))

	models.SelectCourse(5, 1)
	models.SelectCourse(4, 2)
	models.SelectCourse(4, 3)
	models.SelectCourse(1, 3)
	models.SelectCourse(1, 2)
	models.SelectCourse(3, 1)
	models.SelectCourse(3, 2)
	models.SelectCourse(3, 3)
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

	CreateData()
	Router = server.InitRouter()
	// api.Run(":8888")
	m.Run()
}

// 	api.Run(":8888")
// 	os.Exit(m.Run())
