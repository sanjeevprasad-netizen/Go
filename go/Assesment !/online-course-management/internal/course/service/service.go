package service

import (
	"context"
	"gorm.io/gorm"
	courseModel "online-course-management/internal/course/Model"
	"online-course-management/internal/course/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(r *repository.Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(ctx context.Context, c *courseModel.Course) error {
	return s.repo.Create(ctx, c)
}

func (s *Service) GetAll(ctx context.Context) ([]courseModel.Course,error){
	return s.repo.GetAll(ctx);
}

func (s *Service) GetById(ctx context.Context,id uint)(courseModel.Course,error){
	res , err := s.repo.GetById(ctx, id)
	if err == gorm.ErrRecordNotFound 
	{
		return nil, errors.New("course not found")
	}
	return  res , err
}

func (s *Service) Update(ctx context.Context, c *Course,) error {
    return s.repo.Update(ctx, c)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
    return s.repo.Delete(ctx, id)
}