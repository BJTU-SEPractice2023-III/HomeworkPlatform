package models

import (
	"testing"
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

// var (
// 	mock sqlmock.Sqlmock
// 	err  error
// 	db   *sql.DB
// )

// // TestMain是在当前package下，最先运行的一个函数，常用于初始化
// func TestMain(m *testing.M) {
// 	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 	if err != nil {

// 		panic(err)
// 	}
// 	DB, err = gorm.Open("mysql", db)

// 	// m.Run 是调用包下面各个Test函数的入口
// 	os.Exit(m.Run())
// }
