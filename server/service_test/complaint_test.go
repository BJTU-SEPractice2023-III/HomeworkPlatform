package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCreateComplaint(t *testing.T) {
	var cases = []struct {
		Case       string
		HomeworkId uint
		Reason     string
		ExpextCode int
	}{
		{"正确提交", 3, "原因", 200},
		{"未提交作业", 1, "原因", 400},
		{"重复提交", 3, "原因", 400},
		{"空原因", 3, "", 400},
	}
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			// log.Printf("正在测试")
			data := map[string]interface{}{"reason": testcase.Reason}
			jsonData, _ := json.Marshal(data)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/notice/"+strconv.Itoa(int(testcase.HomeworkId)), bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("提交申诉:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestDeleteComplaint(t *testing.T) {
	var cases = []struct {
		Case        string
		ComplaintId uint
		ExpextCode  int
	}{
		{"正确申诉", 1, 200},
		{"不存在这个申诉", 1, 400},
	}
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			// log.Printf("正在测试")
			// data := map[string]interface{}{"reason": testcase.Reason}
			// jsonData, _ := json.Marshal(data)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/v1/notice/"+strconv.Itoa(int(testcase.ComplaintId)), nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("提交申诉:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestGetComplaint(t *testing.T) {
	var cases = []struct {
		Case       string
		HomeworkId uint
		ExpextCode int
	}{
		{"正确查找", 1, 200},
		{"作业没有申诉", 4, 400},
	}
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			// log.Printf("正在测试")
			// data := map[string]interface{}{"reason": testcase.Reason}
			// jsonData, _ := json.Marshal(data)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/notice/"+strconv.Itoa(int(testcase.HomeworkId)), nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("提交申诉:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestUpdateComplaint(t *testing.T) {
	var cases = []struct {
		Case        string
		ComplaintId uint
		Reason      string
		ExpextCode  int
	}{
		{"正确提交", 2, "原因", 200},
		{"complaint不存在", 993, "原因", 400},
		{"空原因", 3, "", 400},
	}
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			// log.Printf("正在测试")
			data := map[string]interface{}{"reason": testcase.Reason}
			jsonData, _ := json.Marshal(data)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/api/v1/notice/"+strconv.Itoa(int(testcase.ComplaintId)), bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("提交申诉:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestSolveComplaint(t *testing.T) {
	var cases = []struct {
		Case        string
		ComplaintId uint
		ExpextCode  int
	}{
		{"正确解决", 2, 200},
		{"complaint不存在", 993, 400},
	}
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			// log.Printf("正在测试")
			// data := map[string]interface{}{"reason": testcase.Reason}
			// jsonData, _ := json.Marshal(data)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/notice/"+strconv.Itoa(int(testcase.ComplaintId))+"/solve", nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("解决申诉:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}
