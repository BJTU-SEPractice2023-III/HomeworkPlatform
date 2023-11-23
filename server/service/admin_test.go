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

func TestGetUsersService(t *testing.T) {
	var cases = []struct {
		Case       string
		UserId     uint
		ExpextCode int
	}{
		{"admin获取", 1, 200},
		{"非admin", 5, 403},
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
			// data := map[string]interface{}{"signature": testcase.Signature}
			// jsonData, _ := json.Marshal(data)
			if testcase.UserId != 1 {
				data := map[string]interface{}{"username": "xb", "password": "123"}
				jsonData, _ := json.Marshal(data)
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				Router.ServeHTTP(w, req)
				Authorization = GetAuthorziation(w)
			}
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/admin/users", nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获得用户列表%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestUserUpdateService(t *testing.T) {
	var cases = []struct {
		Case       string
		UserId     uint
		Username   string
		Password   string
		ExpextCode int
	}{
		{"admin获取", 1, "xbbb", "12", 200},
		{"用户不存在", 1, "xasdb", "12", 400},
		{"非admin", 5, "xbb", "21", 403},
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
			data := map[string]interface{}{"username": testcase.Username, "password": testcase.Password}
			jsonData, _ := json.Marshal(data)
			if testcase.UserId != 1 {
				data := map[string]interface{}{"username": "xb", "password": "123"}
				jsonData, _ := json.Marshal(data)
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				Router.ServeHTTP(w, req)
				Authorization = GetAuthorziation(w)
			}
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/api/v1/admin/users", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获得用户列表%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestDeleteUserService(t *testing.T) {
	var cases = []struct {
		Case         string
		UserId       uint
		TargetUserId uint
		ExpextCode   int
	}{
		{"admin获取", 1, 9, 200},
		{"用户不存在", 1, 9, 400},
		{"非admin", 5, 2, 403},
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
			// data := map[string]interface{}{"username": testcase.Username, "password": testcase.Password}
			// jsonData, _ := json.Marshal(data)
			if testcase.UserId != 1 {
				data := map[string]interface{}{"username": "xb", "password": "123"}
				jsonData, _ := json.Marshal(data)
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				Router.ServeHTTP(w, req)
				Authorization = GetAuthorziation(w)
			}
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/api/v1/admin/users/"+strconv.Itoa(int(testcase.TargetUserId)), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("Admin删除用户%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}
