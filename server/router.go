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
				{
					// GET    api/admin/users     | Get a list of all users
					users.GET("", service.Handler(&service.GetUsersService{}))
					// POST   api/admin/users     | Create a user
					users.POST("", service.Handler(&service.UserUpdateService{}))
					// DELETE api/admin/users/:id | Delete a user
					users.DELETE(":id", service.HandlerWithBindType(&service.DeleteUserService{}, service.BindUri))
				}
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
				homewrok.POST("assign", service.Handler(&service.AssignHomeworkService{}))  // POST api/homework/assign
				homewrok.POST("homeworklists", service.Handler(&service.HomeworkLists{}))   // POST api/homework/homeworklists
				homewrok.POST("delete", service.Handler(&service.DeleteHomework{}))         // POST api/homework/delete
				homewrok.POST("update", service.Handler(&service.UpdateHomeworkService{}))  // POST api/homework/update
				homewrok.GET("information", service.Handler(&service.HomeworkDetail{}))     // GET api/homework/information
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
