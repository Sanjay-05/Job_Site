package service

import (
	"context"
	"jobsite/models"
	"jobsite/repository"
)

type LookupService interface {
	ListCountries(ctx context.Context) ([]models.Country, error)
	ListJobTypes(ctx context.Context) ([]models.JobType, error)
}

type lookupService struct {
	countryRepo repository.CountryRepository
	jobTypeRepo repository.JobTypeRepository
}

func NewLookupService(countryRepo repository.CountryRepository, jobTypeRepo repository.JobTypeRepository) LookupService {
	return &lookupService{countryRepo: countryRepo, jobTypeRepo: jobTypeRepo}
}

func (s *lookupService) ListCountries(ctx context.Context) ([]models.Country, error) {
	return s.countryRepo.List(ctx)
}

func (s *lookupService) ListJobTypes(ctx context.Context) ([]models.JobType, error) {
	return s.jobTypeRepo.List(ctx)
}
