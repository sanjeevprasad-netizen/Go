package enrollmentModel

import "time"

type Enrollment struct {
	StudentID  uint      `json:"studentId" gorm:"primaryKey"`
	CourseID   uint      `json:"courseId" gorm:"primaryKey"`
	EnrolledAt time.Time `json:"enrolledAt"`
}

type EnrollRequest struct {
	StudentID uint `json:"studentId" binding:"required"`
	CourseID  uint `json:"courseId" binding:"required"`
}

type CourseWithEnrollmentCount struct {
	ID              uint      `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	DurationInHours int       `json:"durationInHours"`
	StartDate       time.Time `json:"startDate"`
	Instructor      string    `json:"instructor"`
	EnrollmentCount int64     `json:"enrollmentCount"`
}

