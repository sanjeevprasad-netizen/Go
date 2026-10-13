package database

import (
	"fmt"
	"time"

	courseModel "online-course-management/internal/course/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		return nil, err
	}

	fmt.Println("Database connected")

	// Auto-migrate tables
	if err := db.AutoMigrate(&courseModel.Course{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}
	fmt.Println("Database auto-migrated successfully")

	// Seed dummy data if table is empty
	if err := seedData(db); err != nil {
		return nil, fmt.Errorf("failed to seed dummy data: %w", err)
	}

	return db, nil
}

func seedData(db *gorm.DB) error {
	var count int64
	if err := db.Model(&courseModel.Course{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		dummyCourses := []courseModel.Course{
			{
				Title:           "Introduction to Golang",
				Description:     "Learn the fundamentals of Go programming, syntax, concurrency, and web APIs.",
				DurationInHours: 20,
				StartDate:       time.Now(),
				Instructor:      "Alice Johnson",
			},
			{
				Title:           "Building REST APIs with Gin & GORM",
				Description:     "Hands-on guide to building production-ready RESTful microservices with PostgreSQL.",
				DurationInHours: 35,
				StartDate:       time.Now().AddDate(0, 0, 7),
				Instructor:      "Bob Martin",
			},
		}

		if err := db.Create(&dummyCourses).Error; err != nil {
			return err
		}
		fmt.Println("Dummy data seeded successfully (2 courses added)")
	}

	return nil
}

