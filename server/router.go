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
					// POST   api/v1/courses        | 创建课程
					courses.POST("", service.Handler(&service.CreateCourse{}))
					// PUT    api/v1/courses        | 更新课程
					courses.PUT("", service.Handler(&service.UpdateCourseDescription{}))
					// DELETE api/v1/courses        | 删除课程
					courses.DELETE("", service.Handler(&service.DeleteCourse{}))

					// GET    api/v1/courses/:id/students | 获取指定 id 课程的所有学生信息
					courses.GET(":id/students", service.HandlerWithBindType(&service.GetCourseStudents{}, service.BindUri))
					// POST   api/v1/courses/:id/students | 为指定 id 课程添加学生&选课
					courses.POST(":id/students", service.Handler(&service.AddCourseStudentService{}))

					// GET    api/v1/courses/:id/homeworks | 获取指定 id 课程的所有作业信息
					courses.GET(":id/homeworks", service.HandlerWithBindType(&service.GetCourseHomeworks{}, service.BindUri))
					// POST   api/v1/courses/:id/students | 为指定 id 课程添加作业
					courses.POST(":id/homeworks", service.HandlerNoBind(&service.CreateCourseHomework{}))
				}

				// api/v1/homeworks
				homeworks := auth.Group("homeworks")
				{
					// GET	api/v1/homeworks/:id				| 获取指定 id 作业的信息
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

				submit := auth.Group("submit")
				{
					// POST api/v1/submit
					submit.POST("", service.Handler(&service.SubmitHomework{}))
				}

				//TODO:这里还有点问题
				file := auth.Group("file")
				{
					// GET api/v1/file/:path	| 获得文件
					file.GET(":path", service.HandlerWithBindType(&service.GetFileService{}, service.BindUri))
				}
			}
		}

		// No login required
		user := api.Group("user")
		{
			user.POST("login", service.Handler(&service.UserLoginService{}))       // POST api/user/login
			user.POST("register", service.Handler(&service.UserRegisterService{})) // POST api/user/register
			user.POST("update", service.Handler(&service.UserselfUpdateService{})) // POST api/user/update
		}

		//TODO:后期可以做一下权限验证不能随意获取作业
		file := api.Group("file")
		{
			file.GET("getfile", service.Handler(&service.GetFileService{})) // GET api/file/getfile
		}

		// Login required
		auth := api.Group("")
		auth.Use(middlewares.JWTAuth())
		{
			// Admin required
			admin := api.Group("admin")
			admin.Use(middlewares.AdminCheck())
			{
				users := admin.Group("users")
				// GET    api/admin/users     | Get a list of all users
				users.GET("", service.Handler(&service.GetUsersService{}))
				// POST   api/admin/users     | Create a user
				users.POST("", service.Handler(&service.UserUpdateService{}))
				// DELETE api/admin/users/:id | Delete a user
				users.DELETE(":id", service.HandlerWithBindType(&service.DeleteUserService{}, service.BindUri))
			}

			// api/users
			users := auth.Group("users")
			{
				// GET api/users/:id | Get info of a user
				users.GET(":id", service.HandlerWithBindType(&service.GetUserService{}, service.BindUri))
				// GET api/users/:id/courses | Get courses of a user
				users.GET(":id/courses", service.HandlerWithBindType(&service.GetUserCoursesService{}, service.BindUri))
			}

			//homework_submission
			homework_submission := auth.Group("homeworksubmission")
			{
				//把提交和更新封装一起
				homework_submission.POST("submit", service.Handler(&service.SubmitHomework{})) // POST api/homeworksubmission/submit
			}
			//homework
			homewrok := auth.Group("homework")
			{
				homewrok.POST("assign", service.Handler(&service.AssignHomeworkService{})) // POST api/homework/assign
				homewrok.POST("homeworklists", service.Handler(&service.HomeworkLists{}))  // POST api/homework/homeworklists
				homewrok.POST("delete", service.Handler(&service.DeleteHomework{}))        // POST api/homework/delete
				// GET api/homework/:id | Get homework detail
				homewrok.GET(":id", service.HandlerWithBindType(&service.HomeworkDetail{}, service.BindUri))
				homewrok.POST("update", service.Handler(&service.UpdateHomeworkService{})) // POST api/homework/update
				// homewrok.GET("information", service.Handler(&service.HomeworkDetail{}))     // GET api/homework/information
				homewrok.GET("submitlists", service.Handler(&service.SubmitListsService{})) // GET api/homework/submitlists
			}

			//course
			course := auth.Group("course")
			{
				course.GET("", service.Handler(&service.GetCourses{}))
				course.GET(":id", service.HandlerWithBindType(&service.GetCourse{}, service.BindUri))
				course.POST("create", service.Handler(&service.CreateCourse{}))             // POST api/course/create
				course.POST("update", service.Handler(&service.UpdateCourseDescription{}))  // POST api/course/update
				course.POST("delete", service.Handler(&service.DeleteCourse{}))             // POST api/course/delete
				course.GET("userlists", service.Handler(&service.GetCourseStudentLists{}))  // Get api/course/userlists
				course.POST("select", service.Handler(&service.SelectCourseService{}))      // POST api/course/select
				course.GET("teachingcourse", service.Handler(&service.GetTeachingCourse{})) // GET api/course/teachingcourse
				course.GET("learningcourse", service.Handler(&service.GetLearningCourse{})) // GET api/course/learningcourse
			}

			comment := auth.Group("comment")
			{
				comment.GET("lists", service.Handler(&service.GetCommentListsService{})) // GET api/comment/lists
				comment.POST("", service.Handler(&service.CommentService{}))             // POST api/comment
			}

			grade := auth.Group("grade")
			{
				grade.GET("bysubmissionid", service.Handler(&service.GetGradeBySubmissionIDService{}))  // GET api/grade/bysubmissionid
				grade.GET("byhomeworkid", service.Handler(&service.GetGradeListsByHomeworkIDService{})) // GET api/grade/byhomeworkid
				grade.POST("update", service.Handler(&service.UpdateGradeService{}))                    // POST api/grade/update
			}
		}
	}

	return r
}
