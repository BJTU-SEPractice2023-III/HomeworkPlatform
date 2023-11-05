package server

import (
	"homework_platform/internal/bootstrap"
	"homework_platform/server/middlewares"
	"homework_platform/server/service"

	// "flag"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// TODO: Figure these things out
	config := cors.DefaultConfig()
	config.ExposeHeaders = []string{"Authorization"}
	config.AllowCredentials = true
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// FrontendFS
	if bootstrap.Dev {
		log.Println("Dev flag, using frontend reverse proxy to localhost:5173")
		r.Use(middlewares.FrontendReverseProxy())
	} else {
		r.Use(middlewares.Frontend())
	}

	/*
		路由
	*/
	api := r.Group("api")
	api.Use(gin.Logger())
	api.Use(gin.Recovery())
	{
		v1 := api.Group("v1")
		{
			// No login required
			user := v1.Group("user")
			{
				// TODO: Restful?
				// POST api/v1/user/login | 登录获取 jwt
				user.POST("login", service.Handler(&service.UserLoginService{}))
				// POST api/v1/user       | 注册用户
				user.POST("", service.Handler(&service.UserRegisterService{}))
				// PUT  api/v1/user       | 更新用户信息
				user.PUT("", service.Handler(&service.UserselfUpdateService{}))
			}

			// Admin required
			// api/v1/admin
			admin := v1.Group("admin")
			admin.Use(middlewares.AdminCheck())
			{
				// api/v1/admin/users
				users := admin.Group("users")
				{
					// GET    api/v1/admin/users     | 获取所有用户列表
					users.GET("", service.Handler(&service.GetUsersService{}))
					// POST   api/v1/admin/users     | 创建用户
					users.POST("", service.Handler(&service.UserUpdateService{}))
					// DELETE api/v1/admin/users/:id | 删除用户
					users.DELETE(":id", service.HandlerWithBindType(&service.DeleteUserService{}, service.BindUri))
				}
			}

			// Login required
			auth := v1.Group("")
			auth.Use(middlewares.JWTAuth())
			{
				// api/v1/users
				users := auth.Group("users")
				{
					// GET api/users/:id         | 获取指定 id 用户的信息
					users.GET(":id", service.HandlerWithBindType(&service.GetUserService{}, service.BindUri))
					// GET api/users/:id/courses | 获取指定 id 用户的课程列表（教的课以及学的课）
					users.GET(":id/courses", service.HandlerWithBindType(&service.GetUserCoursesService{}, service.BindUri))
				}

				// api/v1/courses
				courses := auth.Group("courses")
				{
					// GET    api/v1/courses       | 获取所有课程信息
					courses.GET("", service.Handler(&service.GetCourses{}))
					// GET    api/v1/courses/:id   | 获取指定 id 课程信息
					courses.GET(":id", service.HandlerWithBindType(&service.GetCourse{}, service.BindUri))
					// POST   api/v1/course        | 创建课程
					courses.POST("", service.Handler(&service.CreateCourse{}))
					// PUT    api/v1/course        | 更新课程
					courses.PUT("", service.Handler(&service.UpdateCourseDescription{}))
					// DELETE api/v1/course        | 删除课程
					courses.DELETE("", service.Handler(&service.DeleteCourse{}))

					// GET    api/v1/courses/:id/students | 获取指定 id 课程的所有学生信息
					courses.GET(":id/students", service.HandlerWithBindType(&service.GetCourseStudents{}, service.BindUri))
					// POST   api/v1/courses/:id/students | 为指定 id 课程添加学生
					courses.POST(":id/students", service.Handler(&service.AddCourseStudentService{}))

					// GET    api/v1/courses/:id/homeworks | 获取指定 id 课程的所有作业信息
					courses.GET(":id/homeworks", service.HandlerWithBindType(&service.GetCourseHomeworks{}, service.BindUri))
					// POST   api/v1/courses/:id/students | 为指定 id 课程添加作业
					// courses.POST(":id/homeworks", service.Handler(&service.AddCourseHomework{}))
				}

				// api/v1/homeworks
				homeworks := auth.Group("homeworks")
				{
					// GET	api/v1/homeworks/:id 				| 获取指定 id 作业的信息
					homeworks.GET(":id", service.HandlerWithBindType(&service.HomeworkDetail{}, service.BindUri))
					// GET	api/v1/homeworks/:id/submitlists  	| 获取指定 id 作业的全部学生提交信息
					homeworks.GET(":id/submitlists", service.HandlerWithBindType(&service.SubmitListsService{}, service.BindUri))
					// POST api/v1/homeworks 					| 发布作业
					homeworks.POST("", service.Handler(&service.AssignHomeworkService{}))
					// Get api/v1/homeworks/:id/homeworklists 	| 得到指定课程的所有作业
					homeworks.GET(":id/homeworklists", service.HandlerWithBindType(&service.HomeworkLists{}, service.BindUri))
					// DELETE api/v1/homeworks					| 删除指定作业
					homeworks.DELETE("", service.Handler(&service.DeleteHomework{}))
					// PUT api/v1/homeworks  					| 更新作业
					homeworks.PUT("", service.Handler(&service.UpdateHomeworkService{}))
				}

				comment := auth.Group("comment")
				{
					// GET api/v1/comment/:id	| 获得本次作业需要批阅的作业id
					comment.GET(":id", service.HandlerWithBindType(&service.GetCommentListsService{}, service.BindUri))
					// POST api/v1/comment 		|评阅请求提交
					comment.POST("", service.Handler(&service.CommentService{}))
				}

				grade := auth.Group("grade")
				{
					// GET api/v1/:id/grade/bysubmissionid 	| 根据作业提交号获得单个成绩
					grade.GET(":id/bysubmissionid", service.HandlerWithBindType(&service.GetGradeBySubmissionIDService{}, service.BindUri))
					// GET api/v1/:id/grade/byhomeworkid  	| 根据作业号获得一次作业的全部成绩
					grade.GET(":id/byhomeworkid", service.HandlerWithBindType(&service.GetGradeListsByHomeworkIDService{}, service.BindUri))
					// PUT api/v1/grade						| 老师修改成绩
					grade.PUT("", service.Handler(&service.UpdateGradeService{}))
				}

			}
		}

	}

	return r
}
