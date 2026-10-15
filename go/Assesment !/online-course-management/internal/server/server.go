package server

import (
	"net/http"
	coursehandler "online-course-management/internal/course/handler"
	enrollmenthandler "online-course-management/internal/enrollment/handler"
	router_ "online-course-management/internal/router"
	studenthandler "online-course-management/internal/student/handler"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func New(
	port string,
	cH *coursehandler.Handler,
	sH *studenthandler.Handler,
	eH *enrollmenthandler.Handler,
) *Server {
	// Call router setup from router.go
	engine := router_.SetupRouter(cH, sH, eH)
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + port,
			Handler:      engine,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
