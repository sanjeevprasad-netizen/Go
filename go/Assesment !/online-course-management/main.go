package main

import (
	"log"
	"online-course-management/internal/app"
)

func main() {

	application, err := app.New()
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}
	if err := application.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
