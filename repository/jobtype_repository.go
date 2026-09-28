package repository

import (
	"context"
	"jobsite/models"

	"gorm.io/gorm"
)

type JobTypeRepository interface {
	List(context.Context) ([]models.JobType, error)
}

type jobTypeRepository struct {
	db *gorm.DB
}

func NewJobTypeRepository(db *gorm.DB) JobTypeRepository {
	return &jobTypeRepository{db: db}
}

func (r *jobTypeRepository) List(ctx context.Context) ([]models.JobType, error) {
	var types []models.JobType
	if err := r.db.WithContext(ctx).Order("type_name asc").Find(&types).Error; err != nil {
		return nil, err
	}
	return types, nil
}
