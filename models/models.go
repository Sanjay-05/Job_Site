package models

import (
	"jobsite/utility"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return
}

type User struct {
	Base
	UserName string       `gorm:"unique;not null" json:"user_name"`
	Password string       `gorm:"not null" json:"-"`
	Email    string       `gorm:"unique;not null" json:"email"`
	Role     utility.Role `gorm:"not null" json:"role"`
	PhoneNo  string       `gorm:"not null" json:"phone_no"`
}

type Job struct {
	Base
	Website     string     `gorm:"not null" json:"website"`
	CompanyName string     `gorm:"index;not null"`
	JobTitle    string     `gorm:"index;not null"`
	Email       string     `gorm:"not null" json:"email"`
	Description string     `gorm:"column:description" json:"description"`
	Locations   []Location `gorm:"many2many:job_locations;" json:"locations"`
	JobTypes    []JobType  `gorm:"many2many:job_job_types;" json:"job_types"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Country struct {
	Base
	CountryName string `gorm:"unique;not null" json:"country_name"`
}

type Location struct {
	Base
	LocationName string    `json:"location_name"`
	CountryID    uuid.UUID `gorm:"type:uuid;not null" json:"country_id"`
	Country      Country   `gorm:"foreignKey:CountryID" json:"country"`
}

type JobType struct {
	Base
	TypeName string `gorm:"unique;not null" json:"type_name"`
}

type JwtCustomClaims struct {
	UserID string       `json:"user_id"`
	Role   utility.Role `json:"role"`
	jwt.RegisteredClaims
}
