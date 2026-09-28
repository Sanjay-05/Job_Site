package repository

import (
	"context"
	"jobsite/models"
	"jobsite/request"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobRepository interface {
	CreateJob(ctx context.Context, job *models.Job) error
	GetJob(ctx context.Context, filter request.JobFilter) ([]models.Job, int64, error)
	UpdateJob(ctx context.Context, job *models.Job) error
	DeleteJob(ctx context.Context, id uuid.UUID) error
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) CreateJob(ctx context.Context, job *models.Job) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *jobRepository) GetJob(ctx context.Context, filter request.JobFilter) ([]models.Job, int64, error) {
	var jobs []models.Job
	var total int64

	tx := r.db.WithContext(ctx).
		Model(&models.Job{}).
		Joins("LEFT JOIN job_locations jl ON jl.job_id = jobs.id").
		Joins("LEFT JOIN locations l ON l.id = jl.location_id").
		Joins("LEFT JOIN countries c ON c.id = l.country_id").
		Joins("LEFT JOIN job_job_types jjt ON jjt.job_id = jobs.id").
		Joins("LEFT JOIN job_types jt ON jt.id = jjt.job_type_id")

	if filter.Country != "" {
		tx = tx.Where("c.country_name ILIKE ?", "%"+filter.Country+"%")
	}

	if filter.ExcludeCountry != "" {
		tx = tx.Where("c.country_name <> ?", filter.ExcludeCountry)
	}

	if filter.JobType != "" {
		tx = tx.Where("jt.type_name ILIKE ?", "%"+filter.JobType+"%")
	}

	if filter.JobTitle != "" {
		tx = tx.Where("jobs.job_title ILIKE ?", "%"+filter.JobTitle+"%")
	}

	if filter.CompanyName != "" {
		tx = tx.Where("jobs.company_name ILIKE ?", "%"+filter.CompanyName+"%")
	}
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		tx = tx.Where(
			r.db.Where("jobs.job_title ILIKE ?", keyword).
				Or("jobs.description ILIKE ?", keyword))
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		tx = tx.Where(
			r.db.Where("jobs.job_title ILIKE ?", search).
				Or("jobs.description ILIKE ?", search).
				Or("jobs.company_name ILIKE ?", search).
				Or("l.location_name ILIKE ?", search).
				Or("c.country_name ILIKE ?", search))
	}

	if filter.City != "" {
		tx = tx.Where("l.location_name ILIKE ?", "%"+filter.City+"%")
	}

	if err := tx.Session(&gorm.Session{}).
		Distinct("jobs.id").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	offset := (filter.Page - 1) * filter.PageSize

	var jobIDs []uuid.UUID

	idQuery := tx.Session(&gorm.Session{}).
		Select("jobs.id").
		Group("jobs.id").
		Order("MAX(jobs.created_at) DESC").
		Limit(filter.PageSize).
		Offset(offset)

	if err := idQuery.Pluck("jobs.id", &jobIDs).Error; err != nil {
		return nil, 0, err
	}

	if len(jobIDs) == 0 {
		return []models.Job{}, total, nil
	}

	if err := r.db.WithContext(ctx).
		Preload("Locations.Country").
		Preload("JobTypes").
		Where("id IN ?", jobIDs).
		Order("created_at DESC").
		Find(&jobs).Error; err != nil {
		return nil, 0, err
	}
	return jobs, total, nil
}

func (r *jobRepository) UpdateJob(ctx context.Context, job *models.Job) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Model(&models.Job{}).
			Where("id = ?", job.ID).
			Updates(map[string]interface{}{
				"company_name": job.CompanyName,
				"website":      job.Website,
				"job_title":    job.JobTitle,
				"email":        job.Email,
				"description":  job.Description,
			}).Error; err != nil {
			return err
		}

		if err := tx.Model(job).Association("Locations").Replace(job.Locations); err != nil {
			return err
		}

		if err := tx.Model(job).Association("JobTypes").Replace(job.JobTypes); err != nil {
			return err
		}

		return nil
	})
}

func (r *jobRepository) DeleteJob(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var job models.Job
		if err := tx.First(&job, "id = ?", id).Error; err != nil {
			return err
		}

		if err := tx.Model(&job).Association("Locations").Clear(); err != nil {
			return err
		}
		if err := tx.Model(&job).Association("JobTypes").Clear(); err != nil {
			return err
		}

		return tx.Delete(&job).Error
	})
}
