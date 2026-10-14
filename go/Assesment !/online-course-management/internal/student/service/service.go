package service

import (
	"context"
	"errors"
	studentModel "online-course-management/internal/student/model"
	"online-course-management/internal/student/repository"

	"gorm.io/gorm"
)

type Service struct {
	repo *repository.Repository
}

func NewService(r *repository.Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(ctx context.Context, student *studentModel.Student) error {
	return s.repo.Create(ctx, student)
}

func (s *Service) GetAll(ctx context.Context) ([]studentModel.Student, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetById(ctx context.Context, id uint) (studentModel.Student, error) {
	res, err := s.repo.GetById(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return studentModel.Student{}, errors.New("student not found")
	}
	if err != nil {
		return studentModel.Student{}, err
	}
	return *res, nil
}

func (s *Service) Update(ctx context.Context, student *studentModel.Student) error {
	return s.repo.Update(ctx, student)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
