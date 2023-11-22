package service_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCommentService(t *testing.T) {
	var cases = []struct {
		Case       string
		FilePath   string
		ExpextCode int
	}{
		{"正确获得", "data/homework_submission/1/1biey2uhu0g8uc3iioyrcfofo.png.png", 200},
		{"文件不存在", "data/homework_submission/1/1biey2uhu0g8uc3iioyrcfofo.png.png123", 400},
		{"错误的访问", "/homework_submission/1/1biey2uhu0g8uc3iioyrcfofo.png.png123", 400},
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
			req, _ := http.NewRequest("GET", "/api/v1/file/"+testcase.FilePath, nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获得文件:%s,需要的code为%d,但是实际code为%d", testcase.FilePath, testcase.ExpextCode, w.Code)
			}
		})
	}
}
