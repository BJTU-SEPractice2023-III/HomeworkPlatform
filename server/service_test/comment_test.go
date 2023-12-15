package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetGradeBySubmissionIDService(t *testing.T) {
	var cases = []struct {
		Case       string
		HomeworkId uint
		ExpextCode int
	}{
		{"学生正确获得", 5, 200},
		{"老师正确获得", 5, 200},
		{"作业号不存在", 1999, 400},
		{"未开始评阅", 1, 400},
	}
	//登录拿到json
	var Authorization string
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			// log.Printf("正在测试")
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/homeworks/"+strconv.Itoa(int(testcase.HomeworkId))+"/comments", nil)
			// req.Header.Set("Content-Type", "application/json")
			if testcase.Case == "老师正确获得" {
				authorization := GetAuthorziation("xeh", "123")
				req.Header.Set("Authorization", authorization)
			} else {
				req.Header.Set("Authorization", Authorization)
			}

			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获得评论列表:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestGetMyCommentService(t *testing.T) {
	var cases = []struct {
		Case       string
		HomeworkId uint
		ExpextCode int
	}{
		{"正确获得", 5, 200},
		{"作业号不存在", 1999, 400},
		{"未开始评阅", 1, 400},
	}
	//登录拿到json

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/homeworks/5/comments", nil)
	// req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", Authorization)
	Router.ServeHTTP(w, req)
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			// log.Printf("正在测试")
			// data := map[string]interface{}{"name": testcase.Name, "begindate": testcase.BeginDate, "enddate": testcase.EndDate, "description": testcase.Description}
			// jsonData, _ := json.Marshal(data)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/homeworks/"+strconv.Itoa(int(testcase.HomeworkId))+"/mycomments", nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获得评论列表:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestCommentService(t *testing.T) {
	var cases = []struct {
		Case                 string
		HomeworkSubmissionId uint
		Comment              string
		Score                int
		ExpextCode           int
	}{
		{"正确评阅", 4, "kksk", 99, 200},
		{"空内容", 4, "", 99, 400},
		{"小于0的分", 4, "kksk", -1, 400},
		{"大于100的分", 4, "kksk", 101, 400},
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/homeworks/5/comments", nil)
	// req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", Authorization)
	Router.ServeHTTP(w, req)
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			// log.Printf("正在测试")
			data := map[string]interface{}{"score": testcase.Score, "comment": testcase.Comment}
			jsonData, _ := json.Marshal(data)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/comment/"+strconv.Itoa(int(testcase.HomeworkSubmissionId)), bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获得评论列表:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}
