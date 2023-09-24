package server

import (
	"homework_platform/internal/bootstrap"
	"homework_platform/server/middlewares"
	"homework_platform/server/service"

	// "flag"
	// "log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// TODO: Figure these things out
	config := cors.DefaultConfig()
	config.ExposeHeaders = []string{"Authorization"}
	config.AllowCredentials = true
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// FrontendFS
	r.Use(middlewares.Frontend(bootstrap.StaticFS))

	/*
		路由
	*/
	api := r.Group("api")
	{
		// No login required
		user := api.Group("user")
		{
			user.POST("login", service.Handler(&service.UserLoginService{}))       // POST api/user/login
			user.POST("register", service.Handler(&service.UserRegisterService{})) // POST api/user/register
			user.POST("update", service.Handler(&service.UserselfUpdateService{})) // POST api/user/update
		}

		// Login required
		auth := api.Group("")
		auth.Use(middlewares.JWTAuth())
		{
			admin := api.Group("admin")
			admin.Use(middlewares.AdminCheck())
			{
				user := admin.Group("user")
				{
					user.GET("", service.Handler(&service.GetUsersService{}))         // GET  api/admin/user
					user.POST("", service.Handler(&service.UserUpdateService{}))      // POST api/admin/user
					user.POST("delete", service.Handler(&service.DelteUserService{})) //POST api/admin/user/delete
				}
			}

		}

		//homework
		homewrok := api.Group("homework")
		homewrok.Use(middlewares.JWTAuth())
		{
			homewrok.POST("assign", service.Handler(&service.AssignHomeworkService{})) // POST api/homework/assign
			homewrok.POST("homeworklists", service.Handler(&service.HomeworkLists{}))  // POST api/homework/homeworklists
			homewrok.POST("delete", service.Handler(&service.DeleteHomework{}))        // POST api/homework/delete
			homewrok.POST("update", service.Handler(&service.UpdateHomeworkService{})) // POST api/homework/update
			homewrok.GET("information", service.Handler(&service.HomeworkDetail{}))    // GET api/homework/information
		}

		//course
		course := api.Group("course")
		course.Use(middlewares.JWTAuth())
		{
			course.POST("create", service.Handler(&service.CreateCourse{}))             // POST api/course/create
			course.POST("update", service.Handler(&service.UpdateCourseDescription{}))  // POST api/course/update
			course.POST("delete", service.Handler(&service.DeleteCourse{}))             // POST api/course/delete
			course.GET("userlists", service.Handler(&service.GetCourseStudentLists{}))  // Get api/course/userlists
			course.POST("select", service.Handler(&service.SelectCourseService{}))      // POST api/course/select
			course.GET("teachingcourse", service.Handler(&service.GetTeachingCourse{})) // GET api/course/teachingcourse
			course.GET("learningcourse", service.Handler(&service.GetLearningCourse{})) // GET api/course/learningcourse
		}
	}

	return r
}
