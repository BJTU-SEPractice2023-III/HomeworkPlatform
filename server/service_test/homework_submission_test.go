package service_test

import (
	"bytes"
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
		{"成功修改", 3, "1123", 200},
		{"作业不存在", 1, "1123", 400},
		{"作业不存在", 555, "1123", 400},
		{"未选课", 1, "1123", 400},
	}
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
					// fmt.Println(errFile1)
					return
				}
			}
			_ = writer.WriteField("content", testcase.Content)
			err := writer.Close()
			if err != nil {
				// fmt.Println(err)
				return
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/api/v1/submit/"+strconv.Itoa(int(testcase.HomeworkID)), payload)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			if testcase.Case == "未选课" {
				old := Authorization
				Authorization = GetAuthorziation("xbbb", "123")
				req.Header.Set("Authorization", Authorization)
				Authorization = old
			} else {
				req.Header.Set("Authorization", Authorization)
			}
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("提交作业%s,需要的code为%d,但实际为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}
