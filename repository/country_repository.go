package repository

import (
	"context"
	"jobsite/models"

	"gorm.io/gorm"
)

type CountryRepository interface {
	List(ctx context.Context) ([]models.Country, error)
}

type countryRepository struct {
	db *gorm.DB
}

func NewCountryRepository(db *gorm.DB) CountryRepository {
	return &countryRepository{db: db}
}

func (r *countryRepository) List(ctx context.Context) ([]models.Country, error) {
	var countrires []models.Country
	if err := r.db.WithContext(ctx).Order("country_name asc").Find(&countrires).Error; err != nil {
		return nil, err
	}
	return countrires, nil
}
