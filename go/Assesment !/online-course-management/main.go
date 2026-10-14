package main

import (
	"log"
	"net/http"
	"online-course-management/internal/database"
	coursehandler "online-course-management/internal/course/handler"
	courserepo "online-course-management/internal/course/repository"
	courseservice "online-course-management/internal/course/service"
	enrollmenthandler "online-course-management/internal/enrollment/handler"
	enrollmentrepo "online-course-management/internal/enrollment/repository"
	enrollmentservice "online-course-management/internal/enrollment/service"
	studenthandler "online-course-management/internal/student/handler"
	studentrepo "online-course-management/internal/student/repository"
	studentservice "online-course-management/internal/student/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// dsn := os.Getenv("DATABASE_DSN")
	// if dsn == "" {
	// 	log.Fatal("DATABASE_DSN environment variable is required")
	// }
	dsn := "host=localhost user=postgres password=1234 dbname=postgres port=5432 sslmode=disable"
	db, err := database.Connect(dsn)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Course Setup
	cRepo := courserepo.NewRepository(db)
	cService := courseservice.NewService(cRepo)
	cHandler := coursehandler.NewHandler(cService)

	// Student Setup
	sRepo := studentrepo.NewRepository(db)
	sService := studentservice.NewService(sRepo)
	sHandler := studenthandler.NewHandler(sService)

	// Enrollment Setup
	eRepo := enrollmentrepo.NewRepository(db)
	eService := enrollmentservice.NewService(eRepo)
	eHandler := enrollmenthandler.NewHandler(eService)

	router := gin.Default()

	// Course Routes
	router.POST("/course/create", cHandler.Create)
	router.GET("/course", cHandler.GetAll)
	router.GET("/course/:id", cHandler.GetById)
	router.PUT("/courses/:id", cHandler.Update)
	router.DELETE("/courses/:id", cHandler.Delete)

	// Student Routes
	router.POST("/students", sHandler.Create)
	router.GET("/students", sHandler.GetAll)
	router.GET("/students/unenrolled", eHandler.GetUnenrolledStudents)
	router.GET("/students/:id", sHandler.GetById)
	router.PUT("/students/:id", sHandler.Update)
	router.DELETE("/students/:id", sHandler.Delete)

	// Enrollment Routes
	router.POST("/course/enroll", eHandler.Enroll)
	router.GET("/course/:id/students", eHandler.GetStudentsByCourse)
	router.GET("/students/:id/courses", eHandler.GetCoursesByStudent)

	// Custom Query Routes
	router.GET("/courses/top-5", eHandler.GetTopCourses)

	err_r := router.Run(":8080")
	if err_r != nil {
		log.Fatal("Failed to start server:", err_r)
	}

	_ = http.StatusOK
	_ = (*gorm.DB)(nil)
}


