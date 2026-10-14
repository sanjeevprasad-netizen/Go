package repository

import (
	"context"
	courseModel "online-course-management/internal/course/model"
	enrollmentModel "online-course-management/internal/enrollment/model"
	studentModel "online-course-management/internal/student/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// StudentExists checks if student exists
func (r *Repository) StudentExists(ctx context.Context, studentID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&studentModel.Student{}).Where("id = ?", studentID).Count(&count).Error
	return count > 0, err
}

// CourseExists checks if course exists
func (r *Repository) CourseExists(ctx context.Context, courseID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&courseModel.Course{}).Where("id = ?", courseID).Count(&count).Error
	return count > 0, err
}

// IsEnrolled checks if student is already enrolled in the course
func (r *Repository) IsEnrolled(ctx context.Context, studentID, courseID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&enrollmentModel.Enrollment{}).
		Where("student_id = ? AND course_id = ?", studentID, courseID).
		Count(&count).Error
	return count > 0, err
}

// GetStudentEnrollmentCount returns count of active courses for student
func (r *Repository) GetStudentEnrollmentCount(ctx context.Context, studentID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&enrollmentModel.Enrollment{}).
		Where("student_id = ?", studentID).
		Count(&count).Error
	return count, err
}

// GetCourseEnrollmentCount returns count of students in course
func (r *Repository) GetCourseEnrollmentCount(ctx context.Context, courseID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&enrollmentModel.Enrollment{}).
		Where("course_id = ?", courseID).
		Count(&count).Error
	return count, err
}

// Enroll creates an enrollment record
func (r *Repository) Enroll(ctx context.Context, e *enrollmentModel.Enrollment) error {
	return r.db.WithContext(ctx).Create(e).Error
}

// GetCoursesByStudentID returns all courses a student is enrolled in
func (r *Repository) GetCoursesByStudentID(ctx context.Context, studentID uint) ([]courseModel.Course, error) {
	var courses []courseModel.Course
	err := r.db.WithContext(ctx).
		Table("courses").
		Joins("JOIN enrollments ON enrollments.course_id = courses.id").
		Where("enrollments.student_id = ?", studentID).
		Find(&courses).Error
	return courses, err
}

// GetStudentsByCourseID returns all students enrolled in a course
func (r *Repository) GetStudentsByCourseID(ctx context.Context, courseID uint) ([]studentModel.Student, error) {
	var students []studentModel.Student
	err := r.db.WithContext(ctx).
		Table("students").
		Joins("JOIN enrollments ON enrollments.student_id = students.id").
		Where("enrollments.course_id = ?", courseID).
		Find(&students).Error
	return students, err
}

// GetTopEnrolledCourses returns top N courses with the highest enrollment
func (r *Repository) GetTopEnrolledCourses(ctx context.Context, limit int) ([]enrollmentModel.CourseWithEnrollmentCount, error) {
	var results []enrollmentModel.CourseWithEnrollmentCount
	err := r.db.WithContext(ctx).
		Table("courses").
		Select("courses.id, courses.title, courses.description, courses.duration_in_hours, courses.start_date, courses.instructor, COUNT(enrollments.student_id) as enrollment_count").
		Joins("LEFT JOIN enrollments ON enrollments.course_id = courses.id").
		Group("courses.id").
		Order("enrollment_count DESC, courses.id ASC").
		Limit(limit).
		Scan(&results).Error
	return results, err
}

// GetUnenrolledStudents returns students who are not enrolled in any course
func (r *Repository) GetUnenrolledStudents(ctx context.Context) ([]studentModel.Student, error) {
	var students []studentModel.Student
	err := r.db.WithContext(ctx).
		Table("students").
		Joins("LEFT JOIN enrollments ON enrollments.student_id = students.id").
		Where("enrollments.student_id IS NULL").
		Find(&students).Error
	return students, err
}

