package service_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"
)

func TestRegister(t *testing.T) {
	var cases = []struct {
		username   string
		password   string
		ExpextCode int
	}{
		{"1", "2", 200},
		{"1", "2", 400},
		{"", "233", 400},
		{"2", "233", 200},
	}
	url := "http://127.0.0.1:8888/api/v1/user"
	for _, testcase := range cases {
		log.Printf("正在测试")
		data := map[string]interface{}{"username": testcase.username, "password": testcase.password}
		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Fatalf("创建用户testcast转json %s,%s失败", testcase.username, testcase.password)
		}
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		// 获取状态码
		statusCode := resp.StatusCode
		if statusCode != testcase.ExpextCode {
			t.Fatalf("创建用户testcast转json %s,%s,需要的code为%d,但是实际code为%d", testcase.username, testcase.password, testcase.ExpextCode, resp.StatusCode)
		}

		// 获取响应结果
		// body, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	panic(err)
		// }
		// response := string(body)
		// fmt.Println("Response:", response)
		defer resp.Body.Close()
	}
}
