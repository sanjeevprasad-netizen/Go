package database

import (
	"fmt"
	"time"

	courseModel "online-course-management/internal/course/model"
	enrollmentModel "online-course-management/internal/enrollment/model"
	studentModel "online-course-management/internal/student/model"

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
	if err := db.AutoMigrate(&courseModel.Course{}, &studentModel.Student{}, &enrollmentModel.Enrollment{}); err != nil {
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

	var studentCount int64
	if err := db.Model(&studentModel.Student{}).Count(&studentCount).Error; err != nil {
		return err
	}

	if studentCount == 0 {
		dummyStudents := []studentModel.Student{
			{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			{
				Name:  "Jane Smith",
				Email: "jane.smith@example.com",
			},
		}

		if err := db.Create(&dummyStudents).Error; err != nil {
			return err
		}
		fmt.Println("Dummy student data seeded successfully (2 students added)")
	}

	return nil
}
