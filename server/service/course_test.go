package service_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCourse(t *testing.T) {
	var cases = []struct {
		Case       string
		Signature  string
		ExpextCode int
	}{
		{"修改成功", "1", 200},
		{"错误失败", "", 200},
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
			data := map[string]interface{}{"signature": testcase.Signature}
			jsonData, _ := json.Marshal(data)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/api/v1/users/signature", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("修改用户签名为%s,需要的code为%d,但是实际code为%d", testcase.Signature, testcase.ExpextCode, w.Code)
			}
		})
	}
}
