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
	models.CreateUser("xyh", "123") //1
	models.UpgradeToAdmin(1)
	models.CreateUser("xeh", "123") //2
	models.CreateUser("xsh", "123") //3
	models.CreateUser("tjw", "123") //4

	models.CreateUser("xb", "123")         //5
	models.CreateUser("xbb", "123")        //6
	models.CreateUser("xbbb", "123")       //7
	models.CreateUser("xbbbb", "123")      //8
	models.CreateUser("deleteUser", "123") //9
	models.CreateUser("10", "123")

	models.CreateCourse("c++", time.Now(), time.Now().AddDate(0, 0, 1), "c++", 1)
	models.CreateCourse("c3+", time.Now(), time.Now().AddDate(0, 1, 1), "c++", 2)
	models.CreateCourse("c#", time.Now(), time.Now().AddDate(0, 1, 1), "c++", 2)
	models.CreateCourse("c++", time.Now(), time.Now().AddDate(0, 0, 1), "c++", 1)
	models.CreateCourse("c++", time.Now(), time.Now().AddDate(0, 0, 1), "c++", 1)

	//保证作业1是可以提交的
	homework1_id, _ := models.CreateHomework(1, "c++1", "lkksk", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2))
	homework1, _ := models.GetHomeworkByID(homework1_id)
	_ = homework1
	homework2_id, _ := models.CreateHomework(1, "c++2", "lkksk", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2))
	homework2, _ := models.GetHomeworkByID(homework2_id)
	_ = homework2
	homework3_id, _ := models.CreateHomework(2, "c++3", "lkksk", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2))
	homework3, _ := models.GetHomeworkByID(homework3_id)
	homework4_id, _ := models.CreateHomework(2, "c++4", "lkksk", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2))
	homework4, _ := models.GetHomeworkByID(homework4_id)
	homework5_id, _ := models.CreateHomework(2, "c++4", "lkksk", time.Now().AddDate(0, 0, -5), time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2))
	homework5, _ := models.GetHomeworkByID(homework5_id)

	models.SelectCourse(10, 2)
	models.SelectCourse(8, 1)
	models.SelectCourse(6, 1)
	models.SelectCourse(4, 2)
	models.SelectCourse(4, 3)
	models.SelectCourse(1, 3)
	models.SelectCourse(1, 2)
	models.SelectCourse(3, 1)
	models.SelectCourse(3, 2)
	models.SelectCourse(3, 3)

	homework3.AddSubmission(models.HomeworkSubmission{
		UserID:     1,
		HomeworkID: 3,
		Content:    "kksk",
	})
	submission := models.GetHomeWorkSubmissionByID(1)
	submission.Score = 20
	models.DB.Save(&submission)

	homework4.AddSubmission(models.HomeworkSubmission{
		UserID:     1,
		HomeworkID: 4,
		Content:    "kksk",
	})

	homework5.AddSubmission(models.HomeworkSubmission{
		UserID:     1,
		HomeworkID: 5,
		Content:    "kksk",
	})
	homework5.AddSubmission(models.HomeworkSubmission{
		UserID:     10,
		HomeworkID: 5,
		Content:    "kksk",
	})

	homework5.AddSubmission(models.HomeworkSubmission{
		UserID:     77,
		HomeworkID: 5,
		Content:    "kksk",
	})

	models.CreateComment(1, 5, 2)

	models.CreateTeacherComplaint(1, 4, 2, "123")
	models.CreateTeacherComplaint(2, 4, 2, "123")
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

	models.Migrate()

	CreateData()
	Router = server.InitRouter()
	// api.Run(":8888")
	m.Run()
}

// 	api.Run(":8888")
// 	os.Exit(m.Run())
