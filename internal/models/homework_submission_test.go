package models

import (
	"testing"
)

func TestAddHomeworkSubmission(t *testing.T) {
	cases := []struct {
		Case                 string
		user_id, homework_id uint
		content              string
		expected             bool
	}{
		{"正确创建", 2, 2, "kksk", true},
		{"用户不存在", 999, 2, "kksk", false},
		{"课程不存在", 2, 888, "kksk", false},
		{"空内容", 2, 2, "", true},
	}
	t.Helper()
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			homework_submission := HomeworkSubmission{
				UserID:     c.user_id,
				HomeworkID: c.homework_id,
				Content:    c.content,
			}
			res := AddHomeworkSubmission(&homework_submission)
			if res != c.expected {
				t.Fatalf("创建作业提交结果与预期不符!")
			}
		})
	}
}

func TestUpdateself(t *testing.T) {
	homework_submission := GetHomeWorkSubmissionByID(1)
	homework_submission.Content = "kksksss"
	res := homework_submission.UpdateSelf()
	if res != nil {
		t.Fatalf("修改作业提交失败")
	}
}

func TestFindHomeWorkSubmissionByHomeworkIDAndUserID(t *testing.T) {
	cases := []struct {
		Case                 string
		user_id, homework_id uint
		expected             bool
	}{
		{"正确查找", 2, 2, true},
		{"用户号错误", 999, 2, false},
		{"课程号错误", 2, 888, false},
	}
	t.Helper()
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			res := FindHomeWorkSubmissionByHomeworkIDAndUserID(c.homework_id, c.user_id)
			if res == nil && c.expected {
				t.Fatalf("应该找得到但是查找失败")
			} else if res != nil && !c.expected {
				t.Fatalf("应该找得到但是查找失败")
			}
		})
	}
}

func TestGetHomeWorkSubmissionByID(t *testing.T) {
	t.Run("正确查找", func(t *testing.T) {
		res := GetHomeWorkSubmissionByID(1)
		if res == nil {
			t.Fatalf("应该能找到但是失败")
		}
	})

	t.Run("查找失败", func(t *testing.T) {
		res := GetHomeWorkSubmissionByID(999)
		if res != nil {
			t.Fatalf("应该不能找到但是找到了")
		}
	})
}

func TestGetSubmissionListsByHomeworkID(t *testing.T) {
	t.Run("课程号存在", func(t *testing.T) {
		_, err := GetHomeworkByIDWithSubmissionLists(1)
		if err != nil {
			t.Fatalf("应该能找到但是失败")
		}
	})

	t.Run("课程号不存在", func(t *testing.T) {
		_, err := GetHomeworkByIDWithSubmissionLists(99)
		if err == nil {
			t.Fatalf("应该不能找到但是找到了")
		}
	})
}
