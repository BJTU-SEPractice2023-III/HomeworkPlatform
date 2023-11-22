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
)

func TestUpdateSubmission(t *testing.T) {
	var cases = []struct {
		Case       string
		HomeworkID uint
		Content    string
		ExpextCode int
	}{
		{"成功修改", 1, "1123", 200},
		{"作业不存在", 1, "1123", 400},
		{"作业不存在", 555, "1123", 400},
		{"未选课", 1, "1123", 400},
	}
	//登录拿到json
	data := map[string]interface{}{"username": "xbbbb", "password": "123"}
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
