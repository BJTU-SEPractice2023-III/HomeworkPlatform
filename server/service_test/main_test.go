package service_test

import (
	"bytes"
	"encoding/json"
	"homework_platform/internal/bootstrap"
	"homework_platform/internal/config"
	"homework_platform/internal/models"
	"homework_platform/server"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var Router *gin.Engine

type TestCase struct {
	name string
	authorization string
	data any
	code int
}

type ResponseData struct {
	Token string `json:"token"`
}

// 解码 JSON 响应
var responseData struct {
	Data ResponseData `json:"data"`
}

var Authorization string

func GetAuthorziation(username string, password string) string {
	data := map[string]interface{}{"username": username, "password": password}
	w := requestJson("POST", "/api/v1/user/login", data)

	responseBody := w.Body.Bytes()
	err := json.Unmarshal(responseBody, &responseData)
	if err != nil {
		panic(err)
	}

	return responseData.Data.Token
}

func CreateData() {
	user1, _ := models.CreateUser("xyh", "123") // 1
	models.UpgradeToAdmin(1)
	user2, _ := models.CreateUser("xeh", "123") // 2
	user3, _ := models.CreateUser("xsh", "123") // 3
	user4, _ := models.CreateUser("tjw", "123") // 4

	user5, _ := models.CreateUser("xb", "123") // 5
	_ = user5
	user6, _ := models.CreateUser("xbb", "123")  // 6
	user7, _ := models.CreateUser("xbbb", "123") // 7
	_ = user7
	user8, _ := models.CreateUser("xbbbb", "123")      // 8
	user9, _ := models.CreateUser("deleteUser", "123") // 9
	_ = user9
	user10, _ := models.CreateUser("10", "123") // 10

	course1, _ := user1.CreateCourse("c++", time.Now(), time.Now().AddDate(0, 0, 1), "c++")
	course2, _ := user2.CreateCourse("c3+", time.Now(), time.Now().AddDate(0, 1, 1), "c++")
	course3, _ := user2.CreateCourse("c#", time.Now(), time.Now().AddDate(0, 1, 1), "c++")
	course4, _ := user1.CreateCourse("c++", time.Now(), time.Now().AddDate(0, 0, 1), "c++")
	_ = course4
	course5, _ := user1.CreateCourse("c++", time.Now(), time.Now().AddDate(0, 0, 1), "c++")
	_ = course5

	//保证作业1是可以提交的
	homework1, _ := course1.CreateHomework("c++1", "lkksk", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2))
	_ = homework1
	homework2, _ := course1.CreateHomework("c++2", "lkksk", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2))
	_ = homework2
	homework3, _ := course2.CreateHomework("c++3", "lkksk", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2))
	homework4, _ := course2.CreateHomework("c++4", "lkksk", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2))
	homework5, _ := course2.CreateHomework("c++4", "lkksk", time.Now().AddDate(0, 0, -5), time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2))

	user10.SelectCourse(course2.ID)
	user8.SelectCourse(course1.ID)
	user6.SelectCourse(course1.ID)
	user4.SelectCourse(course2.ID)
	user4.SelectCourse(course3.ID)
	user1.SelectCourse(course3.ID)
	user1.SelectCourse(course2.ID)
	user3.SelectCourse(course1.ID)
	user3.SelectCourse(course2.ID)
	user3.SelectCourse(course3.ID)

	homework3.AddSubmission(user1.ID, "kksk")
	submission, _ := models.GetHomeworkSubmissionById(1)
	submission.Score = 20
	models.DB.Save(&submission)

	homework4.AddSubmission(user1.ID, "kksk")

	homework5.AddSubmission(user1.ID, "kksk")
	homework5.AddSubmission(user10.ID, "kksk")

	models.CreateComment(1, 5, 2)

	models.CreateTeacherComplaint(1, 4, 2, "123")
	models.CreateTeacherComplaint(2, 4, 2, "123")
}

func requestJson(method string, url string, data any) *httptest.ResponseRecorder {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic("josn.Marshal error")
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", Authorization)
	Router.ServeHTTP(w, req)

	return w
}

func testRequestWithTestCases(t *testing.T, method string, url string, cases []TestCase) {
	for _, testCase := range cases {
		testRequestWithTestCase(t, method, url, testCase)
	}
}

func testRequestWithTestCase(t *testing.T, method string, url string, testCase TestCase) {
	t.Run(testCase.name, func(t *testing.T) {
		if len(testCase.authorization) != 0 {
			old := Authorization
			Authorization = testCase.authorization
			w := requestJson(method, url, testCase.data)
			assert.Equal(t, w.Code, testCase.code)
			Authorization = old
		} else {
			w := requestJson(method, url, testCase.data)
			assert.Equal(t, w.Code, testCase.code)
		}
	})
}

func TestMain(m *testing.M) {
	bootstrap.Test = true
	models.InitDB()
	bootstrap.Config = &config.Config{JWTSigningString: "moorxJ", SQLDSN: "123"}

	models.Migrate()

	CreateData()
	Router = server.InitRouter()

	Authorization = GetAuthorziation("xyh", "123")
	// api.Run(":8888")
	m.Run()
}

// 	api.Run(":8888")
// 	os.Exit(m.Run())
