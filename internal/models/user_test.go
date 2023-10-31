package models

import (
	"homework_platform/internal/bootstrap"
	"os"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateData() {
	CreateUser("test1", "123")
	CreateUser("test2", "321")
	CreateUser("test3", "kksk")
	CreateUser("test4", "123")
	CreateUser("test5", "123")
	CreateCourse("C++", time.Now(), time.Now().Add(time.Hour), "哈哈", 3)
	CreateCourse("C++", time.Now(), time.Now().Add(time.Hour), "哈哈", 3)
	CreateCourse("C++", time.Now(), time.Now().Add(time.Hour), "哈哈", 3)
	CreateCourse("C++", time.Now(), time.Now().Add(time.Hour), "哈哈", 3)
	CreateCourse("C++", time.Now(), time.Now().Add(time.Hour), "哈哈", 3)
	course1, _ := GetCourseByID(1)
	course2, _ := GetCourseByID(2)
	course1.SelectCourse(2)
	course1.SelectCourse(5)
	course2.SelectCourse(2)
	CreateHomework(2, "原神元素测试", "kksk", time.Now(), time.Now().Add(time.Hour), time.Now().Add(time.Hour).Add(time.Hour))
	CreateHomework(2, "原神元素测试", "kksk", time.Now(), time.Now().Add(time.Hour), time.Now().Add(time.Hour).Add(time.Hour))
	CreateHomework(3, "原神元素测试", "kksk", time.Now(), time.Now().Add(time.Hour), time.Now().Add(time.Hour).Add(time.Hour))
	CreateHomework(1, "原神元素测试", "kksk", time.Now(), time.Now().Add(time.Hour), time.Now().Add(time.Hour).Add(time.Hour))
}

func TestMain(m *testing.M) {
	// 使用 SQLite 内存数据库创建 Gorm 的数据库连接
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	bootstrap.Sqlite = true
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Course{})
	DB.AutoMigrate(&Homework{})
	DB.AutoMigrate(&HomeworkSubmission{})
	DB.AutoMigrate(&Comment{})
	CreateData()
	// 调用包下面各个 Test 函数
	os.Exit(m.Run())
}

// 用户密码测试
func TestCheckPassword(t *testing.T) {
	cases := []struct {
		Case             string
		Name             string
		Target, Password string
		expected         bool
	}{
		{"错误密码", "Admin", "EpK6fkcbLQAihcD3:823e4bebf5915ba903d4bda434457e40d0dc789e", "gzofLRnA2", false},
		{"空的输入密码", "Admin", "EpK6fkcbLQAihcD3:823e4bebf5915ba903d4bda434457e40d0dc789e", "", false},
		{"空用户名", "", "EpK6fkcbLQAihcD3:823e4bebf5915ba903d4bda434457e40d0dc789e", "", false},
		{"正确登录", "Admin", "EpK6fkcbLQAihcD3:823e4bebf5915ba903d4bda434457e40d0dc789e", "gzofLRnA", true},
	}
	t.Helper()
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			user := User{Username: c.Name, Password: c.Target}
			var ans bool
			if ans = user.CheckPassword(c.Password); ans != c.expected {
				t.Fatalf("%s with password %s expected %t, but %t got",
					c.Name, c.Password, c.expected, ans)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	cases := []struct {
		Case     string
		Name     string
		Password string
		expected bool
	}{
		{"正确创建", "xyh", "kksk", true},
		{"空名称", "", "kksk", false},
		{"空密码", "xeh", "", false},
		{"重复注册", "xyh", "kkskkksk", false},
	}
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			t.Helper()
			_, err := CreateUser(c.Name, c.Password)
			if err != nil && c.expected {
				t.Fatalf("create user %s with password %s expected %t, but %s",
					c.Name, c.Password, c.expected, err.Error())
			} else if err == nil && !c.expected {
				t.Fatalf("create user %s with password %s expected %t, but passed",
					c.Name, c.Password, c.expected)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	cases := []struct {
		Case     string
		uid      uint
		expected bool
	}{
		{"查询序号小于1", 0, false},
		{"正确查询", 1, true},
		{"查询序号大于当前最大值", 99, false},
	}
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			_, err := GetUserByID(c.uid)
			if err != nil && c.expected {
				t.Fatalf("get user  by id: %d expected %t, but didn't passed",
					c.uid, c.expected)
			} else if err == nil && !c.expected {
				t.Fatalf("get user  by id: %d expected %t, but passed",
					c.uid, c.expected)
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	cases := []struct {
		Case         string
		new_password string
		expected     bool
	}{
		{"空密码", "", false},
		{"正确修改", "3211231234567", true},
	}
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			user, _ := GetUserByID(1)
			res := user.ChangePassword(c.new_password)
			if res != c.expected {
				t.Fatalf("change user password with new password: %s expected %t, but didn't passed",
					c.new_password, c.expected)
			}
		})
	}
}

func TestDeleteSelf(t *testing.T) {
	user, _ := GetUserByID(1)
	res := user.DeleteSelf()
	if !res {
		t.Fatalf("user delete failed")
	}
}

func TestGetUserByUsername(t *testing.T) {
	cases := []struct {
		Case     string
		username string
		expected bool
	}{
		{"空名称", "", false},
		{"名称正确", "test2", true},
		{"名称不存在", "xyhhh", false},
	}
	for _, c := range cases {
		t.Run(c.Case, func(t *testing.T) {
			_, err := GetUserByUsername(c.username)
			if (err != nil && c.expected) || (err == nil && !c.expected) {
				t.Fatalf("get user by username:%s expected %t but failed",
					c.username, c.expected)
			}
		})
	}
}

func TestGetUserLists(t *testing.T) {
	_, err := GetUserList()
	if err != nil {
		t.Fatalf("get user list failed")
	}
}

func TestGetTeachingCourse(t *testing.T) {
	cases := []struct {
		Case     string
		uid      uint
		expected bool
	}{
		{"课程存在", 3, true},
		{"没有交的课程", 2, true},
	}

	t.Run(cases[0].Case, func(t *testing.T) {
		user, _ := GetUserByID(cases[0].uid)
		courses, err := user.GetTeachingCourse()
		if err != nil && len(courses) == 0 {
			t.Fatalf("get user:%s Teaching course expected %t but failed",
				user.Username, cases[0].expected)
		}
	})

	t.Run(cases[1].Case, func(t *testing.T) {
		user, _ := GetUserByID(cases[1].uid)
		courses, err := user.GetTeachingCourse()
		if err != nil && len(courses) == 0 {
			t.Fatalf("get user:%s Teaching course expected %t but failed",
				user.Username, cases[1].expected)
		}
	})

}

func TestGetLearningCourse(t *testing.T) {
	CreateData()
	cases := []struct {
		Case     string
		uid      uint
		expected bool
	}{
		{"课程存在", 2, true},
		{"没有学习的课程", 3, true},
	}

	t.Run(cases[0].Case, func(t *testing.T) {
		user, _ := GetUserByID(cases[0].uid)
		courses, err := user.GetLearningCourse()
		if err != nil && len(courses) != 2 {
			t.Fatalf("get user:%s Teaching course expected %t but failed",
				user.Username, cases[0].expected)
		}
	})

	t.Run(cases[1].Case, func(t *testing.T) {
		user, _ := GetUserByID(cases[1].uid)
		courses, err := user.GetLearningCourse()
		if err != nil && len(courses) == 0 {
			t.Fatalf("get user:%s Teaching course expected %t but failed",
				user.Username, cases[1].expected)
		}
	})
}
