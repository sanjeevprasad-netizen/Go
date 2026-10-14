package repository

import (
	"context"
	courseModel "online-course-management/internal/course/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, c *courseModel.Course) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *Repository) GetAll(ctx context.Context) ([]courseModel.Course, error) {
	var courses []courseModel.Course
	err := r.db.WithContext(ctx).Find(&courses).Error
	return courses, err
}

func (r *Repository) GetById(ctx context.Context, id uint) (*courseModel.Course, error) {
	var c courseModel.Course

	query := "Select* from courses where id=$1"
	err := r.db.ConnPool.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.Title,
		&c.Description,
		&c.DurationInHours,
		&c.StartDate,
		&c.Instructor,
	)

	//err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *Repository) Update(ctx context.Context, c *courseModel.Course) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&courseModel.Course{}, id).Error
}
