package models

import (
	"testing"
	"time"
)

func TestCreateHomework(t *testing.T) {
	cases := []struct {
		Case                               string
		name                               string
		id                                 uint
		desciption                         string
		beginDate, endDate, commentEndDate time.Time
		expected                           bool
	}{
		{"正确创建", "原神元素测试", 2, "kksk", time.Now(), time.Now().Add(time.Hour), time.Now().Add(time.Hour).Add(time.Hour), true},
		{"空名称", "", 2, "kksk", time.Now(), time.Now().Add(time.Hour), time.Now().Add(time.Hour).Add(time.Hour), false},
		{"空描述", "原神元素测试", 2, "", time.Now(), time.Now().Add(time.Hour), time.Now().Add(time.Hour).Add(time.Hour), false},
		{"开始时间晚于结束时间", "原神元素测试", 2, "kksk", time.Now().Add(time.Hour), time.Now(), time.Now().Add(time.Hour).Add(time.Hour), false},
		{"结束时间晚于批阅时间", "原神元素测试", 2, "kksk", time.Now(), time.Now().Add(time.Hour).Add(time.Hour), time.Now().Add(time.Hour), false},
		{"课程不存在", "原神元素测试", 999, "kksk", time.Now(), time.Now().Add(time.Hour), time.Now().Add(time.Hour).Add(time.Hour), false},
	}
	t.Helper()
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			_, err := CreateHomework(c.id, c.name, c.desciption, c.beginDate, c.endDate, c.commentEndDate)
			if c.expected && err != nil {
				t.Fatalf("创建作业应该报错但是通过")
			} else if !c.expected && err == nil {
				t.Fatalf("创建作业应该通过但是报错")
			}
		})
	}
}

func TestUpdateInformation(t *testing.T) {
	cases := []struct {
		Case                               string
		name                               string
		desciption                         string
		beginDate, endDate, commentEndDate time.Time
		expected                           bool
	}{
		{"正确修改", "原神元素测试", "kksk", time.Now(), time.Now().Add(time.Hour), time.Now().Add(time.Hour).Add(time.Hour), true},
		{"空名称", "", "kksk", time.Now(), time.Now().Add(time.Hour), time.Now().Add(time.Hour).Add(time.Hour), false},
		{"空描述", "原神元素测试", "", time.Now(), time.Now().Add(time.Hour), time.Now().Add(time.Hour).Add(time.Hour), false},
		{"开始时间晚于结束时间", "原神元素测试", "kksk", time.Now().Add(time.Hour), time.Now(), time.Now().Add(time.Hour).Add(time.Hour), false},
		{"结束时间晚于批阅时间", "原神元素测试", "kksk", time.Now(), time.Now().Add(time.Hour).Add(time.Hour), time.Now().Add(time.Hour), false},
	}
	homework, _ := GetHomeworkByID(2)
	t.Helper()
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			res := homework.UpdateInformation(c.name, c.desciption, c.beginDate, c.endDate, c.commentEndDate)
			if !c.expected && res {
				t.Fatalf("创建作业应该报错但是通过")
			} else if c.expected && !res {
				t.Fatalf("创建作业应该通过但是报错")
			}
		})
	}
}

func TestDeleteself(t *testing.T) {
	homework, _ := GetHomeworkByID(3)
	err := homework.Deleteself()
	if err != nil {
		t.Fatalf("作业<id=3>删除自己失败")
	}
}

func TestGetHomeworkByID(t *testing.T) {
	t.Run("id存在", func(t *testing.T) {
		_, err := GetHomeworkByID(2)
		if err != nil {
			t.Fatalf("作业存在但是获得失败")
		}
	})
	t.Run("id不存在", func(t *testing.T) {
		_, err := GetHomeworkByID(999)
		if err == nil {
			t.Fatalf("作业不存在但是获得成功")
		}
	})
}

func TestGetHomeworkByIDWithSubmissionLists(t *testing.T) {
	//TODO:这里还测不了捏
	t.Run("id存在", func(t *testing.T) {
		_, err := GetHomeworkByID(2)
		if err != nil {
			t.Fatalf("作业存在但是获得失败")
		}
	})
	t.Run("id不存在", func(t *testing.T) {
		_, err := GetHomeworkByID(999)
		if err == nil {
			t.Fatalf("作业不存在但是获得成功")
		}
	})
}
