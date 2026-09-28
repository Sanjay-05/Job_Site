package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"jobsite/request"
	"jobsite/service"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Endpoints interface {
	RegisterUser() endpoint.Endpoint
	Login() endpoint.Endpoint
	PostJob() endpoint.Endpoint
	ListJobs() endpoint.Endpoint
	ListCountries() endpoint.Endpoint
	ListJobTypes() endpoint.Endpoint
	UpdateJob() endpoint.Endpoint
	DeleteJob() endpoint.Endpoint
	ForgetPassword() endpoint.Endpoint
	ResetPassword() endpoint.Endpoint
}

type EndpointSet struct {
	AuthSvc   service.UserService
	JobSvc    service.JobService
	LookupSvc service.LookupService
}

func NewEndpoints(AuthSvc service.UserService, JobSvc service.JobService, LookupSvc service.LookupService) Endpoints {
	return &EndpointSet{
		AuthSvc:   AuthSvc,
		JobSvc:    JobSvc,
		LookupSvc: LookupSvc,
	}
}

func (e *EndpointSet) RegisterUser() endpoint.Endpoint {
	return func(ctx context.Context, requests any) (any, error) {
		req := requests.(request.RegisterUserInput)

		err := e.AuthSvc.RegisterUser(ctx, req)
		if err != nil {
			return GenericResponse{Err: err}, nil
		}

		return AuthResponse{Message: "Registration Successful"}, nil
	}
}

func (e *EndpointSet) Login() endpoint.Endpoint {
	return func(ctx context.Context, requests any) (any, error) {
		req := requests.(request.LoginInput)

		token, err := e.AuthSvc.Login(ctx, req.UserName, req.Password)
		if err != nil {
			return GenericResponse{Err: err}, nil
		}

		return AuthResponse{
			Message: "Login successful",
			Token:   token,
		}, nil
	}
}

func (e *EndpointSet) ForgetPassword() endpoint.Endpoint {
	return func(ctx context.Context, requests interface{}) (response interface{}, err error) {
		req := requests.(request.ForgetPasswordRequest)

		err = e.AuthSvc.ForgetPassword(ctx, req.UserName)
		if err != nil {
			return GenericResponse{Err: err}, nil
		}
		return AuthResponse{Message: "Otp send Successfull"}, nil
	}
}

func (e *EndpointSet) ResetPassword() endpoint.Endpoint {
	return func(ctx context.Context, requests interface{}) (response interface{}, err error) {
		req := requests.(request.ResetPasswordRequest)

		err = e.AuthSvc.ResetPassword(ctx, &req)
		if err != nil {
			return GenericResponse{Err: err}, nil
		}
		return AuthResponse{Message: "password reset successfull"}, nil
	}
}

func (e *EndpointSet) PostJob() endpoint.Endpoint {
	return func(ctx context.Context, requests any) (any, error) {

		req, ok := requests.(request.PostJobInput)
		if !ok {
			return GenericResponse{Err: errors.New("invalid request type")}, nil
		}

		err := e.JobSvc.PostJob(ctx, req)
		if err != nil {
			return GenericResponse{Err: err}, nil
		}
		return JobResponse{Message: "Job posted successfully"}, nil
	}
}

func (e *EndpointSet) ListJobs() endpoint.Endpoint {
	return func(ctx context.Context, requests any) (any, error) {
		req := requests.(request.JobFilter)

		result, err := e.JobSvc.ListJobs(ctx, req)
		if err != nil {
			return GenericResponse{Err: err}, nil
		}
		return ListJobsResponse{
			Data:  result.Jobs,
			Total: result.Total,
		}, nil
	}
}

func (e *EndpointSet) ListCountries() endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		countries, err := e.LookupSvc.ListCountries(ctx)
		if err != nil {
			return GenericResponse{Err: err}, nil
		}
		return map[string]any{"data": countries}, nil
	}
}

func (e *EndpointSet) ListJobTypes() endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		types, err := e.LookupSvc.ListJobTypes(ctx)
		if err != nil {
			return GenericResponse{Err: err}, nil
		}
		return map[string]any{"data": types}, nil
	}
}

func (e *EndpointSet) UpdateJob() endpoint.Endpoint {
	return func(ctx context.Context, requests any) (any, error) {
		req := requests.(request.UpdateJobInput)

		err := e.JobSvc.UpdateJob(ctx, req)
		if err != nil {
			return GenericResponse{Err: err}, nil
		}

		return JobResponse{Message: "Job updated successfully"}, nil
	}
}

func (e *EndpointSet) DeleteJob() endpoint.Endpoint {
	return func(ctx context.Context, requests any) (any, error) {
		req := requests.(request.DeleteJobInput)

		err := e.JobSvc.DeleteJob(ctx, req.ID)
		if err != nil {
			return GenericResponse{Err: err}, nil
		}

		return JobResponse{Message: "Job deleted successfully"}, nil
	}
}

type errorer interface {
	error() error
}

type GenericResponse struct {
	Err error `json:"-"`
}

func (r GenericResponse) error() error {
	return r.Err
}

type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Err     error  `json:"-"`
}

func (r AuthResponse) error() error {
	return r.Err
}

type JobResponse struct {
	Message any   `json:"message,omitempty"`
	Err     error `json:"-"`
}

func (r JobResponse) error() error {
	return r.Err
}

type ListJobsResponse struct {
	Data  any   `json:"data,omitempty"`
	Total int64 `json:"total"`
	Err   error `json:"-"`
}

func (r ListJobsResponse) error() error {
	return r.Err
}

func EncodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	EncodeError(context.Background(), err, w)
}

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	switch err {
	case service.ErrValidation:
		code = http.StatusBadRequest
	case service.ErrInvalidCredentials:
		code = http.StatusUnauthorized
	case gorm.ErrRecordNotFound:
		code = http.StatusNotFound
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		EncodeError(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func DecodeRegisterUserRequest(_ context.Context, r *http.Request) (any, error) {
	var req request.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeLoginRequest(_ context.Context, r *http.Request) (any, error) {
	var req request.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeForgetPasswordRequest(_ context.Context, r *http.Request) (any, error) {
	var req request.ForgetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeResetPasswordRequest(_ context.Context, r *http.Request) (any, error) {
	var req request.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodePostJobRequest(_ context.Context, r *http.Request) (any, error) {
	var req request.PostJobInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeListJobsRequest(_ context.Context, r *http.Request) (any, error) {
	q := r.URL.Query()

	page := 1
	pageSize := 10

	if p := q.Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	if ps := q.Get("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	req := request.JobFilter{
		Country:        q.Get("country"),
		ExcludeCountry: q.Get("exclude_country"),
		JobType:        q.Get("job_type"),
		JobTitle:       q.Get("job_title"),
		Keyword:        q.Get("keyword"),
		Search:         q.Get("search"),
		CompanyName:    q.Get("company_name"),
		City:           q.Get("city"),
		Page:           page,
		PageSize:       pageSize,
	}
	return req, nil
}

func DecodeEmptyRequest(_ context.Context, _ *http.Request) (any, error) {
	return struct{}{}, nil
}

func DecodeUpdateJobRequest(_ context.Context, r *http.Request) (any, error) {
	var req request.UpdateJobInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeDeleteJobRequest(_ context.Context, r *http.Request) (any, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid job id:%w", err)
	}
	return request.DeleteJobInput{ID: id}, nil
}
