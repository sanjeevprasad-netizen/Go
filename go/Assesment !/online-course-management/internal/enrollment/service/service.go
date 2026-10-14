package service

import (
	"context"
	"errors"
	"time"

	courseModel "online-course-management/internal/course/model"
	enrollmentModel "online-course-management/internal/enrollment/model"
	"online-course-management/internal/enrollment/repository"
	studentModel "online-course-management/internal/student/model"
)

var (
	ErrStudentNotFound   = errors.New("student not found")
	ErrCourseNotFound    = errors.New("course not found")
	ErrAlreadyEnrolled   = errors.New("student is already enrolled in this course")
	ErrStudentLimitMax   = errors.New("a student cannot enroll in more than 3 courses at once")
	ErrCourseCapacityMax = errors.New("a course cannot have more than 10 students enrolled")
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Enroll(ctx context.Context, studentID, courseID uint) (*enrollmentModel.Enrollment, error) {
	// 1. Check if student exists
	studentExists, err := s.repo.StudentExists(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if !studentExists {
		return nil, ErrStudentNotFound
	}

	// 2. Check if course exists
	courseExists, err := s.repo.CourseExists(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if !courseExists {
		return nil, ErrCourseNotFound
	}

	// 3. Check if already enrolled
	alreadyEnrolled, err := s.repo.IsEnrolled(ctx, studentID, courseID)
	if err != nil {
		return nil, err
	}
	if alreadyEnrolled {
		return nil, ErrAlreadyEnrolled
	}

	// 4. Business Rule: Student cannot enroll in more than 3 courses at once
	studentCoursesCount, err := s.repo.GetStudentEnrollmentCount(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if studentCoursesCount >= 3 {
		return nil, ErrStudentLimitMax
	}

	// 5. Business Rule: Course cannot have more than 10 students enrolled
	courseStudentsCount, err := s.repo.GetCourseEnrollmentCount(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if courseStudentsCount >= 10 {
		return nil, ErrCourseCapacityMax
	}

	// 6. Create enrollment record
	enrollment := &enrollmentModel.Enrollment{
		StudentID:  studentID,
		CourseID:   courseID,
		EnrolledAt: time.Now(),
	}

	if err := s.repo.Enroll(ctx, enrollment); err != nil {
		return nil, err
	}

	return enrollment, nil
}

func (s *Service) GetCoursesByStudentID(ctx context.Context, studentID uint) ([]courseModel.Course, error) {
	studentExists, err := s.repo.StudentExists(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if !studentExists {
		return nil, ErrStudentNotFound
	}

	return s.repo.GetCoursesByStudentID(ctx, studentID)
}

func (s *Service) GetStudentsByCourseID(ctx context.Context, courseID uint) ([]studentModel.Student, error) {
	courseExists, err := s.repo.CourseExists(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if !courseExists {
		return nil, ErrCourseNotFound
	}

	return s.repo.GetStudentsByCourseID(ctx, courseID)
}

func (s *Service) GetTopEnrolledCourses(ctx context.Context, limit int) ([]enrollmentModel.CourseWithEnrollmentCount, error) {
	if limit <= 0 {
		limit = 5
	}
	return s.repo.GetTopEnrolledCourses(ctx, limit)
}

func (s *Service) GetUnenrolledStudents(ctx context.Context) ([]studentModel.Student, error) {
	return s.repo.GetUnenrolledStudents(ctx)
}

