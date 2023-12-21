package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestCreateCourse(t *testing.T) {
	var cases = []struct {
		Case        string
		Name        string
		BeginDate   time.Time
		EndDate     time.Time
		Description string
		ExpextCode  int
	}{
		{"成功创建", "c++", time.Now(), time.Now().Add(time.Minute), "c++课程", 200},
		{"空课程名称", "", time.Now(), time.Now().Add(time.Minute), "c++课程", 400},
		{"开始时间晚于结束", "c++", time.Now().Add(time.Minute), time.Now(), "c++课程", 400},
		{"空描述", "c++", time.Now(), time.Now().Add(time.Minute), "", 200},
	}
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")
			data := map[string]interface{}{"name": testcase.Name, "begindate": testcase.BeginDate, "enddate": testcase.EndDate, "description": testcase.Description}
			jsonData, _ := json.Marshal(data)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/courses", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("创建课程name:%s,需要的code为%d,但是实际code为%d", testcase.Name, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestGetCourses(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/courses", nil)
	req.Header.Set("Authorization", Authorization)
	Router.ServeHTTP(w, req)
	// responseBody := w.Body.Bytes()
	// log.Println(string(responseBody))
}

func TestUpdateCourseDescription(t *testing.T) {
	var cases = []struct {
		Case        string
		CourseID    uint
		Description string
		ExpextCode  int
	}{
		{"成功创建", 1, "kksk", 200},
		{"空描述", 1, "", 200},
		{"非自己的课程", 2, "kksk", 400},
		{"课程号不存在", 99, "kksk", 400},
	}

	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")
			data := map[string]interface{}{"description": testcase.Description}
			jsonData, _ := json.Marshal(data)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/api/v1/courses/"+strconv.Itoa(int(testcase.CourseID)), bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("修改课程描述:%s,需要的code为%d,但是实际code为%d", testcase.Description, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestDeleteCourse(t *testing.T) {
	var cases = []struct {
		Case       string
		CourseID   uint
		ExpextCode int
	}{
		{"权限不足", 5, 400},
		{"课程存在", 5, 200},
		{"课程不存在", 992, 400},
	}

	//登录拿到json
	var Authorization string
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/v1/courses/"+strconv.Itoa(int(testcase.CourseID)), nil)
			// req.Header.Set("Content-Type", "application/json")
			if testcase.Case == "权限不足" {
				authorization := GetAuthorziation("xeh", "123")
				req.Header.Set("Authorization", authorization)
			} else {
				req.Header.Set("Authorization", Authorization)
			}
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("删除课程:%d,需要的code为%d,但是实际code为%d", testcase.CourseID, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestGetCourse(t *testing.T) {
	var cases = []struct {
		Case       string
		CourseID   uint
		ExpextCode int
	}{
		{"成功获得", 4, 200},
		{"课程不存在", 992, 400},
	}

	log.Printf("Authorization为:%s", Authorization)
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/courses/"+strconv.Itoa(int(testcase.CourseID)), nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获得课程:%d,需要的code为%d,但是实际code为%d", testcase.CourseID, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestCreateCourseHomework(t *testing.T) {
	var cases = []struct {
		Case           string
		CourseID       uint
		Name           string
		Description    string
		BeginDate      time.Time
		EndDate        time.Time
		CommentEndDate time.Time
		ExpextCode     int
	}{
		{"成功创建", 1, "c++作业1", "1", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2), 200},
		{"非自己的课程", 2, "c++作业1", "1", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2), 400},
		{"课程不存在", 1232, "c++作业1", "1", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2), 400},
		{"空名称", 1, "", "1", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2), 400},
		{"空描述", 1, "1", "", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2), 400},
		{"时间顺序混乱1", 1, "c++", "1", time.Now(), time.Now().AddDate(0, 1, 1), time.Now().AddDate(0, 0, 2), 400},
		{"时间顺序混乱2", 1, "c--", "1", time.Now().AddDate(0, 1, 1), time.Now(), time.Now().AddDate(0, 0, 2), 400},
	}
	log.Printf("Authorization为:%s", Authorization)
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")

			payload := &bytes.Buffer{}
			writer := multipart.NewWriter(payload)
			file, errFile1 := os.Open("/Users/blackcat/Pictures/1biey2uhu0g8uc3iioyrcfofo.png.png")
			defer file.Close()
			part1,
				errFile1 := writer.CreateFormFile("files", filepath.Base("/Users/blackcat/Pictures/1biey2uhu0g8uc3iioyrcfofo.png.png"))
			_, errFile1 = io.Copy(part1, file)
			if errFile1 != nil {
				fmt.Println(errFile1)
				return
			}
			_ = writer.WriteField("courseid", strconv.Itoa(int(testcase.CourseID)))
			_ = writer.WriteField("description", testcase.Description)
			_ = writer.WriteField("name", testcase.Name)
			_ = writer.WriteField("beginDate", testcase.BeginDate.UTC().Format("2006-01-02T15:04:05Z"))
			_ = writer.WriteField("endDate", testcase.EndDate.UTC().Format("2006-01-02T15:04:05Z"))
			_ = writer.WriteField("commentEndDate", testcase.CommentEndDate.UTC().Format("2006-01-02T15:04:05Z"))
			err := writer.Close()
			if err != nil {
				fmt.Println(err)
				return
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/courses/"+strconv.Itoa(int(testcase.CourseID))+"/homeworks", payload)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("创建课程%s,需要的code为%d,但实际为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestGetCourseHomeworks(t *testing.T) {
	var cases = []struct {
		Case       string
		CourseID   uint
		ExpextCode int
	}{
		{"成功获得", 1, 200},
		{"课程不存在", 992, 400},
	}

	log.Printf("Authorization为:%s", Authorization)
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/courses/"+strconv.Itoa(int(testcase.CourseID))+"/homeworks", nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			log.Print(w.Body)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获得课程:%d,需要的code为%d,但是实际code为%d", testcase.CourseID, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestGetCourseStudents(t *testing.T) {
	var cases = []struct {
		Case       string
		CourseID   uint
		ExpextCode int
	}{
		{"成功获得", 1, 200},
		{"课程不存在", 992, 400},
	}

	log.Printf("Authorization为:%s", Authorization)
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/courses/"+strconv.Itoa(int(testcase.CourseID))+"/students", nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			log.Print(w.Body)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获得课程:%d,需要的code为%d,但是实际code为%d", testcase.CourseID, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestAddCourseStudentService(t *testing.T) {
	var cases = []struct {
		Case       string
		CourseID   uint
		ExpextCode int
	}{
		{"成功添加", 1, 200},
		{"重复添加", 1, 400},
	}

	authorization := GetAuthorziation("tjw", "123")
	log.Printf("Authorization为:%s", Authorization)
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/courses/"+strconv.Itoa(int(testcase.CourseID))+"/students", nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authorization)
			Router.ServeHTTP(w, req)
			log.Print(w.Body)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获得课程:%d,需要的code为%d,但是实际code为%d", testcase.CourseID, testcase.ExpextCode, w.Code)
			}
		})
	}
}
