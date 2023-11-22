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

func TestRegister(t *testing.T) {
	var cases = []struct {
		Case       string
		username   string
		password   string
		ExpextCode int
	}{
		{"正确创建", "1", "2", 200},
		{"重复创建", "1", "2", 400},
		{"没密码", "2", "", 400},
		{"没账户名", "", "233", 400},
	}
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")
			data := map[string]interface{}{"username": testcase.username, "password": testcase.password}
			jsonData, err := json.Marshal(data)
			if err != nil {
				t.Fatalf("创建用户testcast转json %s,%s失败", testcase.username, testcase.password)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			Router.ServeHTTP(w, req)

			if w.Code != testcase.ExpextCode {
				t.Fatalf("创建用户testcast转json %s,%s,需要的code为%d,但是实际code为%d", testcase.username, testcase.password, testcase.ExpextCode, w.Code)
			}

			// 获取响应结果
			// body, err := io.ReadAll(resp.Body)
			// if err != nil {
			// 	panic(err)
			// }
			// response := string(body)
			// fmt.Println("Response:", response)
		})
	}
}

func TestLogin(t *testing.T) {
	var cases = []struct {
		Case       string
		username   string
		password   string
		ExpextCode int
	}{
		{"登陆成功", "xyh", "123", 200},
		{"错误密码", "xyh", "233", 400},
		{"没密码", "2", "", 400},
		{"没账户名", "", "233", 400},
	}
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试")
			data := map[string]interface{}{"username": testcase.username, "password": testcase.password}
			jsonData, _ := json.Marshal(data)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			Router.ServeHTTP(w, req)

			if w.Code != testcase.ExpextCode {
				t.Fatalf("创建用户testcast转json %s,%s,需要的code为%d,但是实际code为%d", testcase.username, testcase.password, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestUpdateUserInformation(t *testing.T) {
	var cases = []struct {
		Case        string
		userName    string
		oldPassword string
		newPassword string
		ExpextCode  int
	}{
		{"修改成功", "xeh", "123", "3", 200},
		{"错误密码", "xsh", "22", "3", 400},
		{"没旧密码", "1", "", "3", 400},
		{"没新密码", "xsh", "22", "", 400},
		{"没账户名", "", "3", "2", 400},
	}
	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			log.Printf("正在测试,用户密码为")
			data := map[string]interface{}{"userName": testcase.userName, "oldPassword": testcase.oldPassword, "newPassword": testcase.newPassword}
			jsonData, _ := json.Marshal(data)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/api/v1/user", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			Router.ServeHTTP(w, req)
			// 获取状态码
			if w.Code != testcase.ExpextCode {
				t.Fatalf("修改用户密码用户testcast %s,%s,需要的code为%d,但是实际code为%d", testcase.userName, testcase.newPassword, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestUpdateSignature(t *testing.T) {
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

func TestGetUserCoursesService(t *testing.T) {
	var cases = []struct {
		Case       string
		UserId     uint
		ExpextCode int
	}{
		{"有课程", 1, 200},
		{"无课程", 5, 200},
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

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/users/"+strconv.Itoa(int(testcase.UserId))+"/courses", nil)
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获得用户课程信息为%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

func TestGetUserNotifications(t *testing.T) {
	var cases = []struct {
		Case       string
		UserId     uint
		ExpextCode int
	}{
		{"有通知", 1, 200},
		{"无通知", 5, 200},
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

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/users/"+strconv.Itoa(int(testcase.UserId))+"/notifications", bytes.NewBuffer(jsonData))
			// req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("测试用例为%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}

