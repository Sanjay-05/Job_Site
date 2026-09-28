package service

import (
	"context"
	"fmt"
	"jobsite/helper"
	"jobsite/models"
	"jobsite/repository"
	"jobsite/request"
	"jobsite/response"

	"github.com/go-kit/log"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobService interface {
	PostJob(ctx context.Context, input request.PostJobInput) error
	ListJobs(ctx context.Context, filter request.JobFilter) (*response.JobListResult, error)
	UpdateJob(ctx context.Context, input request.UpdateJobInput) error
	DeleteJob(ctx context.Context, id uuid.UUID) error
}

type jobService struct {
	jobRepo   repository.JobRepository
	db        repository.DBProvider
	validator *validator.Validate
	logger    log.Logger
}

func NewJobService(jobRepo repository.JobRepository, db repository.DBProvider, logger log.Logger) JobService {
	return &jobService{
		jobRepo: jobRepo,
		db:      db,
		logger:  logger,
	}
}

func (s *jobService) PostJob(ctx context.Context, input request.PostJobInput) error {

	if err := helper.ValidateJson(input); err != nil {
		return err
	}

	tx := s.db.WithContext(ctx)
	var finalLocations []models.Location

	for _, locReq := range input.Locations {
		var loc models.Location
		err := tx.FirstOrCreate(&loc, models.Location{
			LocationName: locReq.Name,
			CountryID:    locReq.CountryID,
		}).Error
		if err != nil {
			return fmt.Errorf("failed to sync location %s: %v", locReq.Name, err)
		}
		finalLocations = append(finalLocations, loc)
	}

	var jobTypes []models.JobType
	if err := tx.Where("id IN ?", input.JobTypeIDs).Find(&jobTypes).Error; err != nil {
		return err
	}

	job := &models.Job{
		CompanyName: input.CompanyName,
		Website:     input.Website,
		JobTitle:    input.JobTitle,
		Email:       input.Email,
		Description: input.Description,
		Locations:   finalLocations,
		JobTypes:    jobTypes,
	}

	return s.jobRepo.CreateJob(ctx, job)
}

func (s *jobService) ListJobs(ctx context.Context, filter request.JobFilter) (*response.JobListResult, error) {
	jobs, total, err := s.jobRepo.GetJob(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(jobs) < 1 {
		return nil, gorm.ErrRecordNotFound
	}
	return &response.JobListResult{
		Jobs:  jobs,
		Total: total,
	}, nil
}

func (s *jobService) UpdateJob(ctx context.Context, input request.UpdateJobInput) error {

	if err := helper.ValidateJson(input); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}

	tx := s.db.WithContext(ctx)

	var finalLocations []models.Location
	for _, locReq := range input.Locations {
		var loc models.Location
		err := tx.FirstOrCreate(&loc, models.Location{
			LocationName: locReq.Name,
			CountryID:    locReq.CountryID,
		}).Error
		if err != nil {
			return err
		}
		finalLocations = append(finalLocations, loc)
	}

	var jobTypes []models.JobType
	if err := tx.Where("id IN ?", input.JobTypeIDs).Find(&jobTypes).Error; err != nil {
		return err
	}

	job := &models.Job{
		Base:        models.Base{ID: input.ID},
		CompanyName: input.CompanyName,
		Website:     input.Website,
		JobTitle:    input.JobTitle,
		Email:       input.Email,
		Description: input.Description,
		Locations:   finalLocations,
		JobTypes:    jobTypes,
	}

	return s.jobRepo.UpdateJob(ctx, job)
}

func (s *jobService) DeleteJob(ctx context.Context, id uuid.UUID) error {
	return s.jobRepo.DeleteJob(ctx, id)
}
