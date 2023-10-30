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

func TestDeleself(t *testing.T) {
	course, _ := GetCourseByID(3)
	err := course.Deleteself()
	if err != nil {
		t.Fatalf("删除课程失败啦")
	}
}

func TestGetStudens(t *testing.T) {
	t.Run("课程学生为空", func(t *testing.T) {
		course, _ := GetCourseByID(4)
		student, err := course.GetStudents()
		if err != nil {
			t.Fatalf("课程学生不存在时报错")
		} else if len(student) != 0 {
			t.Fatalf("课程学生应该为0但是不为0")
		}
	})

	t.Run("课程有学生", func(t *testing.T) {
		course, _ := GetCourseByID(2)
		student, err := course.GetStudents()
		if err != nil {
			t.Fatalf("课程有学生时报错")
		} else if len(student) == 0 {
			t.Fatalf("课程学生不为0但是查询为0")
		}
	})
}

func TestGetCourses(t *testing.T) {
	courses, err := GetCourses()
	if err != nil {
		t.Fatalf("获得所有课程时报错")
	} else if len(courses) == 0 {
		t.Fatalf("获得所有的课程时查询到的课程为0个")
	}
}

func TestFindStudents(t *testing.T) {
	t.Run("课程查询学生存在", func(t *testing.T) {
		course, _ := GetCourseByID(2)
		res := course.GetStudentsByID(2)
		if !res {
			t.Fatalf("通过course查询学生失败!")
		}
	})
	t.Run("课程查询学生不存在", func(t *testing.T) {
		course, _ := GetCourseByID(2)
		res := course.GetStudentsByID(99)
		if res {
			t.Fatalf("通过course查询学生失败!")
		}
	})

}
