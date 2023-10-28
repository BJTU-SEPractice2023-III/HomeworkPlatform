package models

import (
	"testing"
	"time"
)

func TestCreateCourse(t *testing.T) {
	cases := []struct {
		Case               string
		name               string
		beginData, endDate time.Time
		description        string
		teacherid          uint
		expected           bool
	}{
		{"正确创建", "原神", time.Now(), time.Now().Add(time.Hour), "原神,启动!", 1, true},
		{"空描述", "原神", time.Now(), time.Now().Add(time.Hour), "", 1, true},
		{"没有课程名", "", time.Now(), time.Now().Add(time.Hour), "原神,启动!", 1, false},
		{"结束时间早于开始时间", "原神", time.Now().Add(time.Hour), time.Now(), "原神,启动!", 1, false},
	}
	t.Helper()
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			_, err := CreateCourse(c.name, c.beginData, c.endDate, c.description, c.teacherid)
			if (err != nil && c.expected) || (err == nil && !c.expected) {
				t.Fatalf("创建课程期望 %t,但是结果相反", c.expected)
			}
		})
	}
}

func TestGetCourseByID(t *testing.T) {
	t.Run("课程id存在", func(t *testing.T) {
		_, err := GetCourseByID(1)
		if err != nil {
			t.Fatalf("课程id为%d的时候应该能找到该课程,但是没找到", 1)
		}
	})

	t.Run("课程id不存在", func(t *testing.T) {
		_, err := GetCourseByID(999)
		if err == nil {
			t.Fatalf("课程id为%d的时候应该无法找到该课程,但是没有爆错", 999)
		}
	})

}

func TestGetStudens(t *testing.T) {

}
