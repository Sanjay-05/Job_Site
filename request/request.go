package request

import "github.com/google/uuid"

type RegisterUserInput struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
	PhoneNo  string `json:"phone_no" validate:"required"`
}

type LoginInput struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type PostJobInput struct {
	CompanyName string            `json:"company_name" validate:"required"`
	Website     string            `json:"website" validate:"required,url"`
	JobTitle    string            `json:"job_title" validate:"required"`
	Email       string            `json:"email" validate:"required,email"`
	Description string            `json:"description" validate:"required"`
	Locations   []LocationRequest `json:"locations"`
	JobTypeIDs  []uuid.UUID       `json:"job_type_ids"`
}

type LocationRequest struct {
	Name      string    `json:"name"`
	CountryID uuid.UUID `json:"country_id"`
}

type JobFilter struct {
	Country        string
	ExcludeCountry string
	JobType        string
	JobTitle       string
	CompanyName    string
	Keyword        string
	Search         string
	City           string
	Page           int
	PageSize       int
}

type UpdateJobInput struct {
	ID          uuid.UUID         `json:"id" validate:"required"`
	CompanyName string            `json:"company_name" validate:"required"`
	Website     string            `json:"website" validate:"required,url"`
	JobTitle    string            `json:"job_title" validate:"required"`
	Email       string            `json:"email" validate:"required,email"`
	Description string            `json:"description" validate:"required"`
	Locations   []LocationRequest `json:"locations"`
	JobTypeIDs  []uuid.UUID       `json:"job_type_ids"`
}

type DeleteJobInput struct {
	ID uuid.UUID
}

type ForgetPasswordRequest struct {
	UserName string `json:"user_name" validate:"required"`
}

type ResetPasswordRequest struct {
	Username           string `json:"user_name" validate:"required"`
	OTP                string `json:"otp" validate:"required"`
	NewPassword        string `json:"new_password" validate:"required"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required"`
}
