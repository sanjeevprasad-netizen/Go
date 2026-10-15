package router

import (
	coursehandler "online-course-management/internal/course/handler"
	enrollmenthandler "online-course-management/internal/enrollment/handler"
	studenthandler "online-course-management/internal/student/handler"

	"github.com/gin-gonic/gin"
)

// setupRouter configures and returns the Gin router
func SetupRouter(
	cH *coursehandler.Handler,
	sH *studenthandler.Handler,
	eH *enrollmenthandler.Handler,
) *gin.Engine {
	router := gin.Default()

	// Course Routes
	courses := router.Group("/courses")
	{
		courses.POST("", cH.Create)
		courses.GET("", cH.GetAll)
		courses.GET("/:id", cH.GetById)
		courses.PUT("/:id", cH.Update)
		courses.DELETE("/:id", cH.Delete)
		courses.GET("/top-5", eH.GetTopCourses)
		courses.GET("/:id/students", eH.GetStudentsByCourse)
	}

	// Student Routes
	students := router.Group("/students")
	{
		students.POST("", sH.Create)
		students.GET("", sH.GetAll)
		students.GET("/unenrolled", eH.GetUnenrolledStudents)
		students.GET("/:id", sH.GetById)
		students.PUT("/:id", sH.Update)
		students.DELETE("/:id", sH.Delete)
		students.GET("/:id/courses", eH.GetCoursesByStudent)
	}

	// Enrollment Routes
	router.POST("/course/enroll", eH.Enroll)

	return router
}
