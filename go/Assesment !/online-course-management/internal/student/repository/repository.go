package repository

import (
	"context"
	studentModel "online-course-management/internal/student/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, s *studentModel.Student) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *Repository) GetAll(ctx context.Context) ([]studentModel.Student, error) {
	var students []studentModel.Student
	err := r.db.WithContext(ctx).Find(&students).Error
	return students, err
}

func (r *Repository) GetById(ctx context.Context, id uint) (*studentModel.Student, error) {
	var s studentModel.Student
	err := r.db.WithContext(ctx).First(&s, id).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) Update(ctx context.Context, s *studentModel.Student) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&studentModel.Student{}, id).Error
}
