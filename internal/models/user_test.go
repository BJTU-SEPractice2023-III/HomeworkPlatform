package models

import (
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 用户密码测试
func TestCheckPassword(t *testing.T) {
	cases := []struct {
		Name             string
		Target, Password string
		expected         bool
	}{
		{"Admin", "EpK6fkcbLQAihcD3:823e4bebf5915ba903d4bda434457e40d0dc789e", "gzofLRnA2", false},
		{"Admin", "EpK6fkcbLQAihcD3:823e4bebf5915ba903d4bda434457e40d0dc789e", "", false},
		{"", "EpK6fkcbLQAihcD3:823e4bebf5915ba903d4bda434457e40d0dc789e", "", false},
		{"Admin", "EpK6fkcbLQAihcD3:823e4bebf5915ba903d4bda434457e40d0dc789e", "gzofLRnA", true},
	}
	t.Helper()
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			user := User{Username: c.Name, Password: c.Target}
			var ans bool
			if ans = user.CheckPassword(c.Password); ans != c.expected {
				t.Fatalf("%s with password %s expected %t, but %t got",
					c.Name, c.Password, c.expected, ans)
			}
		})
	}
}

var (
	err error
)

func TestMain(m *testing.M) {
	// 使用 SQLite 内存数据库创建 Gorm 的数据库连接
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Course{})
	DB.AutoMigrate(&Homework{})
	DB.AutoMigrate(&HomeworkSubmission{})
	DB.AutoMigrate(&Comment{})
	// 调用包下面各个 Test 函数
	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	cases := []struct {
		Name     string
		Password string
		expected bool
	}{
		{"xyh", "kksk", true},
		{"", "kksk", false},
		{"xeh", "", false},
		{"xyh", "kksk", false},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
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
