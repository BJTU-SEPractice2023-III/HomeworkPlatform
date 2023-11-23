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

func TestGetHomework(t *testing.T) {
	var cases = []struct {
		Case       string
		CourseId   uint
		ExpextCode int
	}{
		{"成功查询", 1, 200},
		{"作业不存在", 999, 400},
	}

	//登录拿到json
	data := map[string]interface{}{"username": "xyh", "password": "123"}
	jsonData, _ := json.Marshal(data)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	Authorization := GetAuthorziation(w)
	log.Printf("Authorization为:%s", Authorization)
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")
			// data := map[string]interface{}{"name": testcase.Name, "begindate": testcase.BeginDate, "enddate": testcase.EndDate, "description": testcase.Description}
			// jsonData, _ := json.Marshal(data)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/homeworks/"+strconv.Itoa(int(testcase.CourseId)), nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获取作业:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestAssignHomeworkService(t *testing.T) {
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
		{"非自己的课程", 3, "c++作业1", "1", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2), 400},
	}
	//登录拿到json
	data := map[string]interface{}{"username": "xyh", "password": "123"}
	jsonData, _ := json.Marshal(data)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	Authorization := GetAuthorziation(w)
	log.Printf("Authorization为:%s", Authorization)
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")

			payload := &bytes.Buffer{}
			writer := multipart.NewWriter(payload)
			if testcase.Case != "空描述" {
				file, errFile1 := os.Open("/Users/blackcat/Pictures/1biey2uhu0g8uc3iioyrcfofo.png.png")
				defer file.Close()
				part1,
					errFile1 := writer.CreateFormFile("files", filepath.Base("/Users/blackcat/Pictures/1biey2uhu0g8uc3iioyrcfofo.png.png"))
				_, errFile1 = io.Copy(part1, file)
				if errFile1 != nil {
					fmt.Println(errFile1)
					return
				}
			}
			_ = writer.WriteField("courseId", strconv.Itoa(int(testcase.CourseID)))
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
			req, _ := http.NewRequest("POST", "/api/homework/assign", payload)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("创建作业%s,需要的code为%d,但实际为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestUpdateHomework(t *testing.T) {
	var cases = []struct {
		Case           string
		HomeworkID     uint
		Name           string
		Description    string
		BeginDate      time.Time
		EndDate        time.Time
		CommentEndDate time.Time
		ExpextCode     int
	}{
		{"成功创建", 1, "c++作业1", "1", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2), 200},
		{"非自己的课程", 3, "c++作业1", "1", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2), 400},
		{"课程不存在", 1232, "c++作业1", "1", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2), 400},
		{"空名称", 1, "", "1", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2), 400},
		// {"空描述", 1, "1", "", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2), 400},
		{"时间顺序混乱1", 1, "c++", "1", time.Now(), time.Now().AddDate(0, 1, 1), time.Now().AddDate(0, 0, 2), 400},
		{"时间顺序混乱2", 1, "c--", "1", time.Now().AddDate(0, 1, 1), time.Now(), time.Now().AddDate(0, 0, 2), 400},
	}
	//登录拿到json
	data := map[string]interface{}{"username": "xyh", "password": "123"}
	jsonData, _ := json.Marshal(data)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	Authorization := GetAuthorziation(w)
	log.Printf("Authorization为:%s", Authorization)
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")

			payload := &bytes.Buffer{}
			writer := multipart.NewWriter(payload)
			if testcase.Case != "空描述" {
				file, errFile1 := os.Open("/Users/blackcat/Pictures/1biey2uhu0g8uc3iioyrcfofo.png.png")
				defer file.Close()
				part1,
					errFile1 := writer.CreateFormFile("files", filepath.Base("/Users/blackcat/Pictures/1biey2uhu0g8uc3iioyrcfofo.png.png"))
				_, errFile1 = io.Copy(part1, file)
				if errFile1 != nil {
					fmt.Println(errFile1)
					return
				}
			}
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
			req, _ := http.NewRequest("PUT", "/api/v1/homeworks/"+strconv.Itoa(int(testcase.HomeworkID)), payload)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("修改作业%s,需要的code为%d,但实际为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestDeleteHomework(t *testing.T) {
	var cases = []struct {
		Case       string
		HomeworkId uint
		ExpextCode int
	}{
		{"成功删除", 2, 200},
		{"非自己的作业", 3, 400},
		{"作业不存在", 999, 400},
	}

	//登录拿到json
	data := map[string]interface{}{"username": "xyh", "password": "123"}
	jsonData, _ := json.Marshal(data)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	Authorization := GetAuthorziation(w)
	log.Printf("Authorization为:%s", Authorization)
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")
			// data := map[string]interface{}{"name": testcase.Name, "begindate": testcase.BeginDate, "enddate": testcase.EndDate, "description": testcase.Description}
			// jsonData, _ := json.Marshal(data)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/v1/homeworks/"+strconv.Itoa(int(testcase.HomeworkId)), nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获取作业:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestSubmitHomework(t *testing.T) {
	var cases = []struct {
		Case       string
		HomeworkID uint
		Content    string
		ExpextCode int
	}{
		{"成功提交", 1, "1123", 200},
		{"重复提交", 1, "1123", 400},
		{"作业不存在", 555, "1123", 400},
		{"未选课", 1, "1123", 400},
	}
	//登录拿到json
	data := map[string]interface{}{"username": "xbb", "password": "123"}
	jsonData, _ := json.Marshal(data)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	Authorization := GetAuthorziation(w)
	log.Printf("Authorization为:%s", Authorization)
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")
			if testcase.Case == "未选课" {
				log.Printf("正在切换用户")
				data := map[string]interface{}{"username": "xbbb", "password": "123"}
				jsonData, _ := json.Marshal(data)
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				Router.ServeHTTP(w, req)
				Authorization = GetAuthorziation(w)
			}
			payload := &bytes.Buffer{}
			writer := multipart.NewWriter(payload)
			if testcase.Case != "空描述" {
				file, errFile1 := os.Open("/Users/blackcat/Pictures/1biey2uhu0g8uc3iioyrcfofo.png.png")
				defer file.Close()
				part1,
					errFile1 := writer.CreateFormFile("files", filepath.Base("/Users/blackcat/Pictures/1biey2uhu0g8uc3iioyrcfofo.png.png"))
				_, errFile1 = io.Copy(part1, file)
				if errFile1 != nil {
					fmt.Println(errFile1)
					return
				}
			}
			_ = writer.WriteField("content", testcase.Content)
			err := writer.Close()
			if err != nil {
				fmt.Println(err)
				return
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/homeworks/"+strconv.Itoa(int(testcase.HomeworkID))+"/submits", payload)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("提交作业%s,需要的code为%d,但实际为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestGetHomeworkSubmission(t *testing.T) {
	var cases = []struct {
		Case       string
		HomeworkId uint
		ExpextCode int
	}{
		{"成功获取", 3, 200},
		{"作业不存在", 1, 400},
	}

	//登录拿到json
	data := map[string]interface{}{"username": "xyh", "password": "123"}
	jsonData, _ := json.Marshal(data)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	Authorization := GetAuthorziation(w)
	log.Printf("Authorization为:%s", Authorization)
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/homeworks/"+strconv.Itoa(int(testcase.HomeworkId))+"/submission", nil)
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获取作业:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}
