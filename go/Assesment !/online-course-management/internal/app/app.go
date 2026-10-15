package app

import (
	"log"
	coursehandler "online-course-management/internal/course/handler"
	courserepo "online-course-management/internal/course/repository"
	courseservice "online-course-management/internal/course/service"
	"online-course-management/internal/database"
	enrollmenthandler "online-course-management/internal/enrollment/handler"
	enrollmentrepo "online-course-management/internal/enrollment/repository"
	enrollmentservice "online-course-management/internal/enrollment/service"
	"online-course-management/internal/server"
	studenthandler "online-course-management/internal/student/handler"
	studentrepo "online-course-management/internal/student/repository"
	studentservice "online-course-management/internal/student/service"
)

type App struct {
	server *server.Server
}

func New() (*App, error) {

	// dsn := os.Getenv("DATABASE_DSN")
	// if dsn == "" {
	// 	log.Fatal("DATABASE_DSN environment variable is required")
	// }
	dsn := "host=localhost user=postgres password=1234 dbname=postgres port=5432 sslmode=disable"
	db, err := database.Connect(dsn)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}

	// Course Setup
	cRepo := courserepo.NewRepository(db)
	cService := courseservice.NewService(cRepo)
	cH := coursehandler.NewHandler(cService)

	// Student Setup
	sRepo := studentrepo.NewRepository(db)
	sService := studentservice.NewService(sRepo)
	sH := studenthandler.NewHandler(sService)

	// Enrollment Setup
	eRepo := enrollmentrepo.NewRepository(db)
	eService := enrollmentservice.NewService(eRepo)
	eH := enrollmenthandler.NewHandler(eService)

	srv := server.New("8080", cH, sH, eH)

	return &App{server: srv}, nil

}
func (a *App) Run() error {
	return a.server.Start()
}
