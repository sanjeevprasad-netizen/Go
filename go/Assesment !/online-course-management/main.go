package main

import (
	"log"
	"net/http"
	"online-course-management/internal/course/database"
	coursehandler "online-course-management/internal/course/handler"
	"online-course-management/internal/course/repository"
	"online-course-management/internal/course/service"

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

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	coursehandler := coursehandler.NewHandler(service)

	router := gin.Default()

	router.POST("/course/create", coursehandler.Create)
	router.GET("/course", coursehandler.GetAll)
	router.GET("/course/:id", coursehandler.GetById)
	router.PUT("/courses/:id", coursehandler.Update)
	router.DELETE("/courses/:id", coursehandler.Delete)

	err_r := router.Run(":8080")
	if err_r != nil {
		log.Fatal("Failed to start server:", err_r)
	}

	_ = http.StatusOK
	_ = (*gorm.DB)(nil)

}
