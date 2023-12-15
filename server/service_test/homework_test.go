package service_test

import (
	"bytes"
	"fmt"
	"homework_platform/server/service"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetHomework(t *testing.T) {
	var cases = []TestCase{
		{"成功查询", "", service.GetHomeworkById{ID: 1}, 200},
		{"作业不存在", "", service.GetHomeworkById{ID: 999}, 400},
	}
	for _, testCase := range cases {
		testRequestWithTestCase(
			t,
			"GET",
			fmt.Sprintf("/api/v1/homeworks/%d", testCase.data.(service.GetHomeworkById).ID),
			testCase,
		)
	}
}

// 这个先屎着
func TestAssignHomeworkService(t *testing.T) {
	var cases = []TestCase{
		{"成功创建", "", service.CreateCourseHomework{
			CourseID:       1,
			Name:           "c++作业1",
			Description:    "1",
			BeginDate:      time.Now(),
			EndDate:        time.Now().AddDate(0, 0, 1),
			CommentEndDate: time.Now().AddDate(0, 0, 2),
		}, 200},
		{"非自己的课程", "", service.CreateCourseHomework{
			CourseID:       3,
			Name:           "c++作业1",
			Description:    "1",
			BeginDate:      time.Now(),
			EndDate:        time.Now().AddDate(0, 0, 1),
			CommentEndDate: time.Now().AddDate(0, 0, 2),
		}, 400},
	}
	for _, testcase := range cases {
		t.Run(testcase.name, func(t *testing.T) {

			payload := &bytes.Buffer{}

			writer := multipart.NewWriter(payload)
			os.WriteFile("test_file.txt", []byte("全测了"), 0666)
			file, _ := os.Open("./test_file.txt")
			defer func() {
				file.Close()
				os.Remove("./test_file.txt")
			}()

			part1, _ := writer.CreateFormFile("files", filepath.Base("/Users/blackcat/Pictures/1biey2uhu0g8uc3iioyrcfofo.png.png"))
			io.Copy(part1, file)

			data := testcase.data.(service.CreateCourseHomework)
			_ = writer.WriteField("courseId", strconv.Itoa(int(data.CourseID)))
			_ = writer.WriteField("description", data.Description)
			_ = writer.WriteField("name", data.Name)
			_ = writer.WriteField("beginDate", data.BeginDate.UTC().Format("2006-01-02T15:04:05Z"))
			_ = writer.WriteField("endDate", data.EndDate.UTC().Format("2006-01-02T15:04:05Z"))
			_ = writer.WriteField("commentEndDate", data.CommentEndDate.UTC().Format("2006-01-02T15:04:05Z"))
			err := writer.Close()
			if err != nil {
				fmt.Println(err)
				return
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/courses/%d/homeworks", data.CourseID), payload)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			assert.Equal(t, w.Code, testcase.code)
		})
	}
}

// 这个先彻底屎着，后面更新拆成信息和文件
// func TestUpdateHomework(t *testing.T) {
// 	var cases = []TestCase{
// 		{"成功创建", "", service.UpdateHomework{
// 			HomeworkID: 1,
// 			Name: "c++作业1",
// 			Description: "1",
// 			BeginDate: time.Now(),
// 			EndDate: time.Now().AddDate(0, 0, 1),
// 			CommentEndDate: time.Now().AddDate(0, 0, 2),
// 		}, 200},
// 		{"非自己的课程", "", service.UpdateHomework{
// 			HomeworkID: 3,
// 			Name: "c++作业1",
// 			Description: "1",
// 			BeginDate: time.Now(),
// 			EndDate: time.Now().AddDate(0, 0, 1),
// 			CommentEndDate: time.Now().AddDate(0, 0, 2),
// 		}, 400},
// 		{"课程不存在", "", service.UpdateHomework{
// 			HomeworkID: 1232,
// 			Name: "c++作业1",
// 			Description: "1",
// 			BeginDate: time.Now(),
// 			EndDate: time.Now().AddDate(0, 0, 1),
// 			CommentEndDate: time.Now().AddDate(0, 0, 2),
// 		 }, 400},
// 		{"空名称", "", service.UpdateHomework{
// 			HomeworkID: 1,
// 			Name: "",
// 			Description: "1",
// 			BeginDate: time.Now(),
// 			EndDate: time.Now().AddDate(0, 0, 1),
// 			CommentEndDate: time.Now().AddDate(0, 0, 2),
// 		}, 400},
// 		// {"空描述", 1, "1", "", time.Now(), time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2), 400},
// 		{"时间顺序混乱1", "", service.UpdateHomework{
// 			HomeworkID: 1,
// 			Name: "c++",
// 			Description: "1",
// 			BeginDate: time.Now(),
// 			EndDate: time.Now().AddDate(0, 1, 1),
// 			CommentEndDate: time.Now().AddDate(0, 0, 2),
// 		}, 400},
// 		{"时间顺序混乱2", "", service.UpdateHomework{
// 			HomeworkID: 1,
// 			Name: "c--",
// 			Description: "1",
// 			BeginDate: time.Now().AddDate(0, 1, 1),
// 			EndDate: time.Now(),
// 			CommentEndDate: time.Now().AddDate(0, 0, 2),
// 		}, 400},
// 	}
// 	for _, testcase := range cases {
// 		t.Run(testcase.name, func(t *testing.T) {
// 			// log.Printf("正在测试")

// 			payload := &bytes.Buffer{}
// 			writer := multipart.NewWriter(payload)
// 			// if testcase.Case != "空描述" {
// 			// 	file, errFile1 := os.Open("/Users/blackcat/Pictures/1biey2uhu0g8uc3iioyrcfofo.png.png")
// 			// 	defer file.Close()
// 			// 	part1,
// 			// 		errFile1 := writer.CreateFormFile("files", filepath.Base("/Users/blackcat/Pictures/1biey2uhu0g8uc3iioyrcfofo.png.png"))
// 			// 	_, errFile1 = io.Copy(part1, file)
// 			// 	if errFile1 != nil {
// 			// 		fmt.Println(errFile1)
// 			// 		return
// 			// 	}
// 			// }
// 			_ = writer.WriteField("description", testcase.Description)
// 			_ = writer.WriteField("name", testcase.Name)
// 			_ = writer.WriteField("beginDate", testcase.BeginDate.UTC().Format("2006-01-02T15:04:05Z"))
// 			_ = writer.WriteField("endDate", testcase.EndDate.UTC().Format("2006-01-02T15:04:05Z"))
// 			_ = writer.WriteField("commentEndDate", testcase.CommentEndDate.UTC().Format("2006-01-02T15:04:05Z"))
// 			err := writer.Close()
// 			if err != nil {
// 				fmt.Println(err)
// 				return
// 			}

// 			w := httptest.NewRecorder()
// 			req, _ := http.NewRequest("PUT", "/api/v1/homeworks/"+strconv.Itoa(int(testcase.HomeworkID)), payload)
// 			req.Header.Set("Content-Type", writer.FormDataContentType())
// 			req.Header.Set("Authorization", Authorization)
// 			Router.ServeHTTP(w, req)
// 			if w.Code != testcase.ExpextCode {
// 				t.Fatalf("修改作业%s,需要的code为%d,但实际为%d", testcase.Case, testcase.ExpextCode, w.Code)
// 			}
// 		})
// 	}
// }

func TestDeleteHomeworkById(t *testing.T) {
	var cases = []TestCase{
		{"成功删除", "", service.DeleteHomeworkById{ID: 2}, 200},
		{"非自己的作业", "", service.DeleteHomeworkById{ID: 3}, 400},
		{"作业不存在", "", service.DeleteHomeworkById{ID: 999}, 400},
	}

	for _, testcase := range cases {
		testRequestWithTestCase(
			t,
			"DELETE",
			fmt.Sprintf("/api/v1/homeworks/%d", testcase.data.(service.DeleteHomeworkById).ID),
			testcase,
		)
	}
}

func TestSubmitHomework(t *testing.T) {
	var cases = []struct {
		Case       string
		HomeworkID uint
		Content    string
		ExpextCode int
	}{
		{"成功提交", 1, "1123", 200},
		{"重复提交", 1, "1123", 400},
		{"作业不存在", 555, "1123", 400},
		{"未选课", 1, "1123", 400},
	}

	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			// log.Printf("正在测试")
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
			if testcase.Case == "未选课" {
				oldAuthorization := Authorization
				Authorization = GetAuthorziation("xbbb", "123")
				req.Header.Set("Authorization", Authorization)
				Authorization = oldAuthorization
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

func TestGetHomeworkSubmission(t *testing.T) {
	var cases = []struct {
		Case       string
		HomeworkId uint
		ExpextCode int
	}{
		{"成功获取", 3, 200},
		{"作业不存在", 1, 400},
	}

	for _, testcase := range cases {
		t.Run(testcase.Case, func(t *testing.T) {
			// log.Printf("正在测试")
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/homeworks/"+strconv.Itoa(int(testcase.HomeworkId))+"/submission", nil)
			req.Header.Set("Authorization", Authorization)
			Router.ServeHTTP(w, req)
			if w.Code != testcase.ExpextCode {
				t.Fatalf("获取作业:%s,需要的code为%d,但是实际code为%d", testcase.Case, testcase.ExpextCode, w.Code)
			}
		})
	}
}
