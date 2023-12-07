package server

import (
	"homework_platform/internal/bootstrap"
	"homework_platform/server/middlewares"
	"homework_platform/server/service"
	user_service "homework_platform/server/service/user"
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
		v2 := api.Group("v2")
		{
			v2.GET("file/:id", service.HandlerBindUri(&service.DownloadFileById{}))

			auth := v2.Group("")
			auth.Use(middlewares.JWTAuth())
			{
				auth.GET("notifications", service.HandlerBindUri(&user_service.GetNotifications{}))
			}
		}

		v1 := api.Group("v1")
		{
			// No login required
			user := v1.Group("user")
			{
				// TODO: Restful?
				// POST api/v1/user/login | 登录获取 jwt
				user.POST("login", service.HandlerBind(&user_service.Login{}))
				// POST api/v1/user       | 注册用户
				user.POST("", service.HandlerBind(&user_service.Register{}))
				// PUT api/v1/users       | 更新用户信息
				user.PUT("", service.HandlerBind(&user_service.UserselfupdateService{}))
			}

			// Admin required
			// api/v1/admin
			// admin := v1.Group("admin")
			// admin.Use(middlewares.AdminCheck())
			// {
			// 	// api/v1/admin/users
			// 	users := admin.Group("users")
			// 	{
			// 		// GET    api/v1/admin/users     | 获取所有用户列表
			// 		users.GET("", service.HandlerBind(&service.GetUsersService{}))
			// 		// PUT   api/v1/admin/users     | 修改密码
			// 		users.PUT("", service.HandlerBind(&service.UserUpdateService{}))
			// 		// DELETE api/v1/admin/users/:id | 删除用户
			// 		users.DELETE(":id", service.HandlerWithBindType(&service.DeleteUserService{}, service.BindUri))
			// 	}
			// }

			// Login required
			auth := v1.Group("")
			auth.Use(middlewares.JWTAuth())
			{
				// api/v1/users
				users := auth.Group("users")
				{
					// GET api/v1/users/:id         		| 获取指定 id 用户的信息
					users.GET(":id", service.HandlerBindUri(&user_service.GetUserService{}))
					// GET api/v1/users/name				| 获得用户的姓名
					users.GET("/name",service.HandlerBind(&user_service.GetUserNameService{}))
					// PUT api/v1/users/signature 			| 更新用户的签名
					users.PUT("signature", service.HandlerBind(&user_service.UpdateSignature{}))
					// GET api/v1/users/:id/courses 		| 获取指定 id 用户的课程列表（教的课以及学的课）
					users.GET(":id/courses", service.HandlerBindUri(&user_service.GetUserCourses{}))
					// GET api/v1/users/:id/notifications	| 获得指定 id 用户的提示信息
					// ! deprecated: use GET api/v2/nitifications
					users.GET(":id/notifications", service.HandlerBindUri(&user_service.GetUserNotifications{}))
				}

				// api/v1/courses
				courses := auth.Group("courses")
				{
					// GET    api/v1/courses        | 获取所有课程信息
					courses.GET("", service.HandlerBind(&service.GetCourses{}))
					// GET    api/v1/courses/:id    | 获取指定 id 课程信息
					courses.GET(":id", service.HandlerBindUri(&service.GetCourse{}))
					// POST   api/v1/courses        | 创建课程
					courses.POST("", service.HandlerBind(&service.CreateCourse{}))
					// PUT    api/v1/courses        | 更新课程
					courses.PUT(":id", service.HandlerWithBindType(&service.UpdateCourseDescription{}, service.Bind|service.BindUri))
					// DELETE api/v1/courses/:id        | 删除课程
					courses.DELETE(":id", service.HandlerBindUri(&service.DeleteCourse{}))

					// GET    api/v1/courses/:id/students | 获取指定 id 课程的所有学生信息
					courses.GET(":id/students", service.HandlerBindUri(&service.GetCourseStudents{}))
					// POST   api/v1/courses/:id/students | 为指定 id 课程添加学生（请求提交者）
					courses.POST(":id/students", service.HandlerBindUri(&service.AddCourseStudentService{}))

					// GET    api/v1/courses/:id/homeworks | 获取指定 id 课程的所有作业信息
					courses.GET(":id/homeworks", service.HandlerBindUri(&service.GetCourseHomeworks{}))
					// POST   api/v1/courses/:id/homeworks | 为指定 id 课程添加作业
					courses.POST(":id/homeworks", service.HandlerWithBindType(&service.CreateCourseHomework{}, service.Bind|service.BindUri))
				}

				// api/v1/homeworks
				homeworks := auth.Group("homeworks")
				{
					// GET    api/v1/homeworks/:id               | 获取指定 id 作业的信息
					homeworks.GET(":id", service.HandlerBindUri(&service.GetHomeworkById{}))
					// PUT    api/v1/homeworks/:id               | 更新指定 id 作业的信息
					homeworks.PUT(":id", service.HandlerWithBindType(&service.UpdateHomework{}, service.Bind|service.BindUri))
					// DELETE api/v1/homeworks/:id               | 删除指定 id 作业
					homeworks.DELETE(":id", service.HandlerBindUri(&service.DeleteHomeworkById{}))
					// POST   api/v1/homeworks/:id/submits       | 上传指定 id 作业的提交
					homeworks.POST(":id/submits", service.HandlerWithBindType(&service.SubmitHomework{}, service.Bind|service.BindUri))

					// GET 	  api/v1/homeworks/:id/comments		 | 得到id作业号的用户应该批阅的列表
					homeworks.GET(":id/comments", service.HandlerBindUri(&service.GetCommentListsService{}))
					// GET	  api/v1/homeworks/:id/mycomments 	| 得到id作业号的用户的被评论信息
					homeworks.GET(":id/mycomments", service.HandlerBindUri(&service.GetMyCommentService{}))
					// GET 	  api/v1/homeworks/:id/submission 	|	根据作业id和用户id获取作业信息
					homeworks.GET(":id/submission", service.HandlerBindUri(&service.GetHomeworkUserSubmission{}))
				}

				notice := auth.Group("notice")
				{
					// POST 	api/v1/notice/:id       | 提交对homeworkId的申诉
					notice.POST(":id", service.HandlerWithBindType(&service.CreateComplaint{}, service.Bind|service.BindUri))
					// DELETE 	api/v1/notice/:id       | 撤销指定Id的申诉
					notice.DELETE(":id", service.HandlerBindUri(&service.DeleteComplaint{}))
					// PUT 		api/v1/notice/:id       | 修改指定id的申诉
					notice.PUT(":id", service.HandlerWithBindType(&service.UpdateComplaint{}, service.Bind|service.BindUri))
					// POST 	api/v1/notice/:id/solve | 老师确认申诉
					notice.POST(":id/solve", service.HandlerBindUri(&service.SolveComplaint{}))
					// GET		api/v1/notice/:id       | 得到某次作业提交的申诉,学生得到自己的申诉,老师得到全部的申诉
					notice.GET(":id", service.HandlerBindUri(&service.GetComplaint{}))
				}

				comment := auth.Group("comment")
				{
					// GET api/v1/comment/:id       | 获得作业信息
					// comment.GET(":id", service.HandlerWithBindType(&service.CommentService{}, service.BindUri))
					// POST api/v1/comment/:id      | 评阅请求提交,提交和修改一体化接口()
					comment.POST(":id", service.HandlerWithBindType(&service.CommentService{}, service.Bind|service.BindUri))
				}

				grade := auth.Group("grade")
				{
					// GET api/v1/grade/:id/bysubmissionid      | 根据作业提交号获得单个成绩
					grade.GET(":id/bysubmissionid", service.HandlerBindUri(&service.GetGradeBySubmissionIDService{}))
					// GET api/v1/grade/:id                     | 根据作业号获得指定作业的成绩,其中老师一次获得全部,学生获得自己的
					grade.GET(":id", service.HandlerBindUri(&service.GetGradeListsByHomeworkIDService{}))
					// PUT api/v1/grade/:id	                    | 老师修改成绩
					grade.PUT(":id", service.HandlerWithBindType(&service.UpdateGradeService{}, service.Bind|service.BindUri))
				}

				submit := auth.Group("submit")
				{
					// POST api/v1/submit                       | 提交作业
					submit.POST("", service.HandlerBind(&service.SubmitHomework{}))
					// PUT api/v1/submit/:id                    | 修改作业提交信息
					submit.PUT(":id", service.HandlerWithBindType(&service.UpdateUserSubmission{}, service.Bind|service.BindUri))
					// GET api/v1/submit/:id                    | 获得指定submission_id的作业
					submit.GET(":id", service.HandlerBindUri(&service.GetSubmissionByIdService{}))
				}

				file := auth.Group("file")
				{
					// GET api/v1/file/:path                    | 获得文件
					file.GET("*path", service.HandlerBind(&service.GetFileService{}))
				}

				ai := auth.Group("ai")
				{
					// POST api/v1/ai/gpt	| 询问G哥
					ai.POST("gpt", service.HandlerBind(&service.GPTService{}))
					// POST api/v1/ai/spark | 连接 SSE
					ai.GET("spark", service.HandlerNoBind(&service.ConnectSpark{}))
					// POST api/v1/ai/spark
					ai.POST("spark", service.HandlerBind(&service.SparkService{}))
					// POST api/v1/ai/spark/image
					ai.POST("spark/image", service.HandlerBind(&service.SparkImageService{}))
				}
			}
		}
	}

	return r
}
