package models

import "testing"

func TestCreateComment(t *testing.T) {
	cases := []struct {
		Case                                         string
		homework_sibmission_id, user_id, homework_id uint
		expected                                     bool
	}{
		{"正确创建", 1, 2, 2, true},
		{"作业提交不存在", 999, 1, 2, false},
		{"用户不存在", 2, 888, 2, false},
	}
	t.Helper()
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			res := CreateComment(c.homework_sibmission_id, c.user_id, c.homework_id)
			if res != c.expected {
				t.Fatalf("创建评论结果与预期不符!")
			}
		})
	}
}

func TestUpdateSelf(t *testing.T) {
	comment, _ := GetCommentByUserIDAndHomeworkSubmissionID(3, 1)
	c := comment.(Comment)
	if c.UpdateSelf("kkk", 50) != nil {
		t.Fatalf("更新自己失败!")
	}
}

func TestGetCommentByUserIDAndHomeworkSubmissionID(t *testing.T) {
	_, err := GetCommentByUserIDAndHomeworkSubmissionID(2, 1)
	if err != nil {
		t.Fatalf("根据用户id和作业提交id获取评论失败")
	}
}

func TestGetCommentListsByUserIDAndHomeworkID(t *testing.T) {
	_, err := GetCommentListsByUserIDAndHomeworkID(2, 1)
	if err != nil {
		t.Fatalf("根据用户id和作业提交id获取评论失败")
	}
}

func TestGetCommentBySubmissionID(t *testing.T) {
	t.Run("有评论数目", func(t *testing.T) {
		comment, err := GetCommentBySubmissionID(1)
		if err != nil || len(comment) == 0 {
			t.Fatalf("获得评论失败")
		}
	})
	t.Run("无评论数目", func(t *testing.T) {
		comment, err := GetCommentBySubmissionID(2)
		if err != nil || len(comment) != 0 {
			t.Fatalf("获得评论失败")
		}
	})

}
