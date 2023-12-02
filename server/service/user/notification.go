package user

import (
	"homework_platform/internal/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	TeachingHomeworkToFinishNotification = iota
	TeachingHomeworkToCommentNotification
	ComplaintToBeSolvedNotification
	ComplaintInProgressNotification
	LearningHomeworkToFinishNotification
	LearningHomeworkToCommentNotification
)

type Notification struct {
	NotificationType uint `json:"notificationType"`
	NotificationData any
}

type GetNotifications struct {
	userId uint `uri:"id" binding:"required"`
}

func (service *GetNotifications) Handle(c *gin.Context) (any, error) {
	user, err := models.GetUserByID(service.userId)
	if err != nil {
		return nil, err
	}

	courses, err := user.GetCourses()
	if err != nil {
		return nil, err
	}

	// Learning homework notifications
	for _, learningCourse := range courses.LearningCourses {
		homeworks, err := learningCourse.GetHomeworks()
		if homeworks == nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(homeworks); j++ {
			// 在批阅时段中
			if homeworks[j].CommentEndDate.After(time.Now()) {
				// 作业已经开始
				if homeworks[j].BeginDate.Before(time.Now()) {
					// 作业在提交时段内
					if homeworks[j].EndDate.After(time.Now()) {
						// homework := models.GetHomeWorkSubmissionByHomeworkIDAndUserID(homeworks[j].ID, user.ID)
						// 没交作业
						// if homework == nil {
						// 	notifications.LeaningHomeworkListsToFinish =
						// 		append(notifications.LeaningHomeworkListsToFinish, homeworks[j])
						// }
					} else {
						// 评论时段内,获取所有的comment
						comments, err := models.GetCommentListsByUserIDAndHomeworkID(user.ID, homeworks[j].ID)
						if err != nil {
							return nil, err
						}
						// 如果有score==-1就代表尚未完成评论
						for i := 0; i < len(comments); i++ {
							if comments[i].Score == -1 {
								// notifications.LeaningHomeworkListsToComment =
								// 	append(notifications.TeachingHomeworkListsToComment, homeworks[j])
								// break
							}
						}
					}
				}
			}
		}
	}

	return nil, nil
}

type GetUserNotifications struct {
	ID uint `uri:"id" binding:"required"`
}

type Notifications struct {
	Type string `json:"type"`

	TeachingHomeworkListsToFinish  []models.Homework `json:"homeworkInProgress"`
	TeachingHomeworkListsToComment []models.Homework `json:"commentInProgress"`

	ComplaintToBeSolved []models.Complaint `json:"complaintToBeSolved"`
	ComplaintInProgress []models.Complaint `json:"complaintInProgress"`

	LeaningHomeworkListsToFinish  []models.Homework `json:"homeworksToBeCompleted"`
	LeaningHomeworkListsToComment []models.Homework `json:"commentToBeCompleted"`
}

// 返回应该尚未提交的作业,待批阅的作业和每门课最新发布的作业
func (service *GetUserNotifications) Handle(c *gin.Context) (any, error) {
	user, err := models.GetUserByID(service.ID)
	if err != nil {
		return nil, err
	}
	courses, err := user.GetCourses()
	if err != nil {
		return nil, err
	}
	var notifications Notifications
	//得到教的课中进行中和批阅中的作业
	log.Printf("len of homework%d\n", len(courses.LearningCourses))
	//得到学的课中还没完成的作业和还没批阅的作业
	for _, course := range courses.LearningCourses {
		//每门课的作业
		homeworks, err := course.GetHomeworks()
		if homeworks == nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(homeworks); j++ {
			// 在批阅时段中
			if homeworks[j].CommentEndDate.After(time.Now()) {
				// 作业已经开始
				if homeworks[j].BeginDate.Before(time.Now()) {
					// 作业在提交时段内
					if homeworks[j].EndDate.After(time.Now()) {
						homework := models.GetHomeWorkSubmissionByHomeworkIDAndUserID(homeworks[j].ID, user.ID)
						// 没交作业
						if homework == nil {
							notifications.LeaningHomeworkListsToFinish =
								append(notifications.LeaningHomeworkListsToFinish, homeworks[j])
						}
					} else {
						// 评论时段内,获取所有的comment
						comments, err := models.GetCommentListsByUserIDAndHomeworkID(user.ID, homeworks[j].ID)
						if err != nil {
							return nil, err
						}
						// 如果有score==-1就代表尚未完成评论
						for i := 0; i < len(comments); i++ {
							if comments[i].Score == -1 {
								notifications.LeaningHomeworkListsToComment =
									append(notifications.TeachingHomeworkListsToComment, homeworks[j])
								break
							}
						}
					}
				}
			}
		}
	}
	//得到老师的课正在进行的作业
	for _, course := range courses.TeachingCourses {
		// 教的课中的作业
		homeworks, err := course.GetHomeworks()
		if homeworks == nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(homeworks); j++ {
			// comment尚未结束
			if homeworks[j].CommentEndDate.After(time.Now()) {
				//作业已经开始
				if homeworks[j].BeginDate.Before(time.Now()) {
					//在提交时段内
					if homeworks[j].EndDate.After(time.Now()) {
						notifications.TeachingHomeworkListsToFinish =
							append(notifications.TeachingHomeworkListsToFinish, homeworks[j])
					} else {
						notifications.TeachingHomeworkListsToComment =
							append(notifications.TeachingHomeworkListsToComment, homeworks[j])
					}
				}
			}
		}
	}
	//得到老师待审核的complaint
	notifications.ComplaintToBeSolved, err = models.GetComplaintByTeacherID(user.ID)
	if err != nil {
		return nil, err
	}
	//得到学生还未被处理的complaint
	notifications.ComplaintInProgress, err = models.GetComplaintByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
