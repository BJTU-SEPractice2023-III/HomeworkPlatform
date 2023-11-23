package service_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetCommentListsService(t *testing.T) {
	var cases = []struct {
		Case         string
		SubmissionId uint
		ExpextCode   int
	}{
		{"正确获得", 1, 200},
		{"作业号不存在", 1999, 400},
	}
	//登录拿到json
	var Authorization string
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")

			data := map[string]interface{}{"username": "xeh", "password": "123"}
			jsonData, _ := json.Marshal(data)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			Router.ServeHTTP(w, req)
			Authorization = GetAuthorziation(w)

			w = httptest.NewRecorder()
			req, _ = http.NewRequest("GET", "/api/v1/grade/"+strconv.Itoa(int(testcase.SubmissionId))+"/bysubmissionid", nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获得成绩:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestGetGradeListsByHomeworkIDService(t *testing.T) {
	var cases = []struct {
		Case       string
		HomeworkId uint
		ExpextCode int
	}{
		{"正确获得", 3, 200},
		{"无权限", 1, 400},
		{"作业号不存在", 1999, 400},
	}
	//登录拿到json
	var Authorization string
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")

			data := map[string]interface{}{"username": "xeh", "password": "123"}
			jsonData, _ := json.Marshal(data)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			Router.ServeHTTP(w, req)
			Authorization = GetAuthorziation(w)

			w = httptest.NewRecorder()
			req, _ = http.NewRequest("GET", "/api/v1/grade/"+strconv.Itoa(int(testcase.HomeworkId))+"/byhomeworkid", nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获得成绩:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestUpdateGradeService(t *testing.T) {
	var cases = []struct {
		Case                 string
		HomeworkSubmissionId uint
		Score                int
		ExpextCode           int
	}{
		{"正确修改", 1, 99, 200},
		{"小于0", 3, -1, 400},
		{"大于0", 3, 101, 400},
		{"无权限", 3, 99, 400},
		{"作业号不存在", 1999, 99, 400},
	}
	//登录拿到json
	var Authorization string
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")
			if testcase.Case == "无权限" {
				data := map[string]interface{}{"username": "xyh", "password": "123"}
				jsonData, _ := json.Marshal(data)
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				Router.ServeHTTP(w, req)
				Authorization = GetAuthorziation(w)
			} else {
				data := map[string]interface{}{"username": "xeh", "password": "123"}
				jsonData, _ := json.Marshal(data)
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				Router.ServeHTTP(w, req)
				Authorization = GetAuthorziation(w)
			}
			data := map[string]interface{}{"score": testcase.Score}
			jsonData, _ := json.Marshal(data)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/api/v1/grade/"+strconv.Itoa(int(testcase.HomeworkSubmissionId)), bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获得成绩:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}
