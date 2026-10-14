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

	// Seed dummy data if tables are empty
	if err := seedData(db); err != nil {
		return nil, fmt.Errorf("failed to seed dummy data: %w", err)
	}

	return db, nil
}

func seedData(db *gorm.DB) error {
	// 1. Seed Courses
	var courseCount int64
	if err := db.Model(&courseModel.Course{}).Count(&courseCount).Error; err != nil {
		return err
	}

	if courseCount == 0 {
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
			{
				Title:           "PostgreSQL & Database Design",
				Description:     "Master database normalization, indexing, querying, and performance tuning.",
				DurationInHours: 25,
				StartDate:       time.Now().AddDate(0, 0, 14),
				Instructor:      "Carol White",
			},
		}

		if err := db.Create(&dummyCourses).Error; err != nil {
			return err
		}
		fmt.Println("Dummy courses seeded successfully (3 courses added)")
	}

	// 2. Seed Students
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
			{
				Name:  "Alex Rivera",
				Email: "alex.rivera@example.com",
			},
			{
				Name:  "Sarah Connor",
				Email: "sarah.connor@example.com",
			},
		}

		if err := db.Create(&dummyStudents).Error; err != nil {
			return err
		}
		fmt.Println("Dummy students seeded successfully (4 students added)")
	}

	// 3. Seed Enrollments
	var enrollmentCount int64
	if err := db.Model(&enrollmentModel.Enrollment{}).Count(&enrollmentCount).Error; err != nil {
		return err
	}

	if enrollmentCount == 0 {
		var courses []courseModel.Course
		var students []studentModel.Student

		if err := db.Find(&courses).Error; err != nil {
			return err
		}
		if err := db.Find(&students).Error; err != nil {
			return err
		}

		if len(courses) >= 2 && len(students) >= 4 {
			dummyEnrollments := []enrollmentModel.Enrollment{
				{
					StudentID:  students[0].ID,
					CourseID:   courses[0].ID,
					EnrolledAt: time.Now(),
				},
				{
					StudentID:  students[0].ID,
					CourseID:   courses[1].ID,
					EnrolledAt: time.Now(),
				},
				{
					StudentID:  students[1].ID,
					CourseID:   courses[0].ID,
					EnrolledAt: time.Now(),
				},
				{
					StudentID:  students[2].ID,
					CourseID:   courses[1].ID,
					EnrolledAt: time.Now(),
				},
				// students[3] (Sarah Connor) remains unenrolled to showcase GET /students/unenrolled
			}

			if err := db.Create(&dummyEnrollments).Error; err != nil {
				return err
			}
			fmt.Println("Dummy enrollments seeded successfully (4 enrollments added)")
		}
	}

	return nil
}
