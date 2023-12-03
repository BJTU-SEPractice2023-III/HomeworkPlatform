package service_test

import (
	"fmt"
	service "homework_platform/server/service/user"
	"testing"
)

func TestRegister(t *testing.T) {
	var cases = []TestCase{
		{"正确创建", "", service.Register{Username: "1", Password: "2"}, 200},
		{"重复创建", "", service.Register{Username: "1", Password: "2"}, 400},
		{"没密码", "", service.Register{Username: "", Password: "2"}, 400},
		{"没账户名", "", service.Register{Username: "2", Password: ""}, 400},
	}
	testRequestWithTestCases(t, "POST", "/api/v1/user", cases)
}

func TestLogin(t *testing.T) {
	var cases = []TestCase{
		{"登陆成功", "", service.Login{Username: "xyh", Password: "123"}, 200},
		{"错误密码", "", service.Login{Username: "xyh", Password: "233"}, 400},
		{"没密码", "", service.Login{Username: "2", Password: ""}, 400},
		{"没账户名", "", service.Login{Username: "", Password: "233"}, 400},
	}
	testRequestWithTestCases(t, "POST", "/api/v1/user/login", cases)
}

func TestUpdateUserInformation(t *testing.T) {
	var cases = []TestCase{
		{"修改成功", "", service.UserselfupdateService{UserName: "xeh", OldPassword: "123", NewPassword: "3"}, 200},
		{"错误密码", "", service.UserselfupdateService{UserName: "xsh", OldPassword: "22", NewPassword: "3"}, 400},
		{"没旧密码", "", service.UserselfupdateService{UserName: "1", OldPassword: "", NewPassword: "3"}, 400},
		{"没新密码", "", service.UserselfupdateService{UserName: "xsh", OldPassword: "22", NewPassword: ""}, 400},
		{"没账户名", "", service.UserselfupdateService{UserName: "", OldPassword: "3", NewPassword: "2"}, 400},
	}
	testRequestWithTestCases(t, "PUT", "/api/v1/user", cases)
}

func TestUpdateSignature(t *testing.T) {
	var cases = []TestCase{
		{"修改成功", "", service.UpdateSignature{Signature: "1"}, 200},
		{"错误失败", "", service.UpdateSignature{Signature: ""}, 200},
	}
	testRequestWithTestCases(t, "PUT", "/api/v1/users/signature", cases)
}

func TestGetUserCoursesService(t *testing.T) {
	var cases = []TestCase{
		{"有课程", "", service.GetUserCourses{ID: 1}, 200},
		{"无课程", "", service.GetUserCourses{ID: 5}, 200},
	}
	for _, testCase := range cases {
		testRequestWithTestCase(t, "GET", fmt.Sprintf("/api/v1/users/%d/courses", testCase.data.(service.GetUserCourses).ID), testCase)
	}
}

func TestGetUserNotifications(t *testing.T) {
	var cases = []TestCase{
		{"有通知", GetAuthorziation("xyh", "123"), service.GetNotifications{}, 200},
		{"无通知", GetAuthorziation("xb", "123"), service.GetNotifications{}, 200},
	}

	testRequestWithTestCases(t, "GET", "/api/v2/notifications", cases)
}
