package handler_test

import (
	"context"
	"errors"
	"jobsite/handler"
	"jobsite/models"
	"jobsite/request"
	"jobsite/response"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockJobService struct {
	mock.Mock
}

func (m *MockJobService) PostJob(ctx context.Context, req request.PostJobInput) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockJobService) ListJobs(ctx context.Context, filter request.JobFilter) (*response.JobListResult, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*response.JobListResult), args.Error(1)
}

func (m *MockJobService) UpdateJob(ctx context.Context, req request.UpdateJobInput) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockJobService) DeleteJob(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestJobEndpoints(t *testing.T) {
	locationID := uuid.New()
	jobTypeID := uuid.New()
	jobID := uuid.New()

	jobInput := request.PostJobInput{
		CompanyName: "Google",
		Website:     "https://google.com",
		JobTitle:    "Go Developer",
		Email:       "hr@google.com",
		Description: "Backend developer role",
		Locations: []request.LocationRequest{
			{
				Name:      "Chennai",
				CountryID: locationID,
			},
		},
		JobTypeIDs: []uuid.UUID{jobTypeID},
	}

	updateInput := request.UpdateJobInput{
		ID:          jobID,
		CompanyName: "Google",
		Website:     "https://google.com",
		JobTitle:    "Senior Go Developer",
		Email:       "hr@google.com",
		Description: "Backend role",
		Locations: []request.LocationRequest{
			{
				Name:      "Bangalore",
				CountryID: locationID,
			},
		},
		JobTypeIDs: []uuid.UUID{jobTypeID},
	}

	listFilter := request.JobFilter{Page: 1, PageSize: 10}
	jobList := &response.JobListResult{
		Jobs: []models.Job{
			{JobTitle: "Dev 1", CompanyName: "Google"},
			{JobTitle: "Dev 2", CompanyName: "Amazon"},
		},
		Total: 2,
	}

	tests := []struct {
		name       string
		endpointFn func(svc *MockJobService) func(ctx context.Context, req any) (any, error)
		input      any
		mockSetup  func(svc *MockJobService)
		expectResp any
		expectErr  bool
	}{
		{
			name:       "PostJob success",
			endpointFn: func(svc *MockJobService) func(ctx context.Context, req any) (any, error) { return handler.NewEndpoints(nil, svc, nil).PostJob() },
			input:      jobInput,
			mockSetup: func(svc *MockJobService) {
				svc.On("PostJob", mock.Anything, jobInput).Return(nil)
			},
			expectResp: handler.JobResponse{Message: "Job posted successfully"},
			expectErr:  false,
		},
		{
			name:       "PostJob error",
			endpointFn: func(svc *MockJobService) func(ctx context.Context, req any) (any, error) { return handler.NewEndpoints(nil, svc, nil).PostJob() },
			input:      jobInput,
			mockSetup: func(svc *MockJobService) {
				svc.On("PostJob", mock.Anything, jobInput).Return(errors.New("db error"))
			},
			expectResp: handler.GenericResponse{Err: errors.New("db error")},
			expectErr:  true,
		},
		{
			name:       "PostJob invalid request",
			endpointFn: func(svc *MockJobService) func(ctx context.Context, req any) (any, error) { return handler.NewEndpoints(nil, svc, nil).PostJob() },
			input:      "invalid",
			mockSetup:  func(svc *MockJobService) {},
			expectResp: handler.GenericResponse{Err: errors.New("invalid request type")},
			expectErr:  true,
		},
		{
			name:       "ListJobs success",
			endpointFn: func(svc *MockJobService) func(ctx context.Context, req any) (any, error) { return handler.NewEndpoints(nil, svc, nil).ListJobs() },
			input:      listFilter,
			mockSetup: func(svc *MockJobService) {
				svc.On("ListJobs", mock.Anything, listFilter).Return(jobList, nil)
			},
			expectResp: handler.ListJobsResponse{Data: jobList.Jobs, Total: jobList.Total},
			expectErr:  false,
		},
		{
			name:       "ListJobs error",
			endpointFn: func(svc *MockJobService) func(ctx context.Context, req any) (any, error) { return handler.NewEndpoints(nil, svc, nil).ListJobs() },
			input:      listFilter,
			mockSetup: func(svc *MockJobService) {
				svc.On("ListJobs", mock.Anything, listFilter).Return(&response.JobListResult{}, errors.New("db error"))
			},
			expectResp: handler.GenericResponse{Err: errors.New("db error")},
			expectErr:  true,
		},
		{
			name:       "UpdateJob success",
			endpointFn: func(svc *MockJobService) func(ctx context.Context, req any) (any, error) { return handler.NewEndpoints(nil, svc, nil).UpdateJob() },
			input:      updateInput,
			mockSetup: func(svc *MockJobService) {
				svc.On("UpdateJob", mock.Anything, updateInput).Return(nil)
			},
			expectResp: handler.JobResponse{Message: "Job updated successfully"},
			expectErr:  false,
		},
		{
			name:       "UpdateJob error",
			endpointFn: func(svc *MockJobService) func(ctx context.Context, req any) (any, error) { return handler.NewEndpoints(nil, svc, nil).UpdateJob() },
			input:      updateInput,
			mockSetup: func(svc *MockJobService) {
				svc.On("UpdateJob", mock.Anything, updateInput).Return(errors.New("update failed"))
			},
			expectResp: handler.GenericResponse{Err: errors.New("update failed")},
			expectErr:  true,
		},
		{
			name:       "DeleteJob success",
			endpointFn: func(svc *MockJobService) func(ctx context.Context, req any) (any, error) { return handler.NewEndpoints(nil, svc, nil).DeleteJob() },
			input:      request.DeleteJobInput{ID: jobID},
			mockSetup: func(svc *MockJobService) {
				svc.On("DeleteJob", mock.Anything, jobID).Return(nil)
			},
			expectResp: handler.JobResponse{Message: "Job deleted successfully"},
			expectErr:  false,
		},
		{
			name:       "DeleteJob error",
			endpointFn: func(svc *MockJobService) func(ctx context.Context, req any) (any, error) { return handler.NewEndpoints(nil, svc, nil).DeleteJob() },
			input:      request.DeleteJobInput{ID: jobID},
			mockSetup: func(svc *MockJobService) {
				svc.On("DeleteJob", mock.Anything, jobID).Return(errors.New("delete failed"))
			},
			expectResp: handler.GenericResponse{Err: errors.New("delete failed")},
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockJobService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockSvc)
			}

			ep := tt.endpointFn(mockSvc)
			resp, err := ep(context.Background(), tt.input)

			assert.NoError(t, err)

			if tt.expectErr {
				r := resp.(handler.GenericResponse)
				assert.Error(t, r.Err)
			} else {
				if !reflect.DeepEqual(resp, tt.expectResp) {
					t.Errorf("expected %+v, got %+v", tt.expectResp, resp)
				}
			}
			mockSvc.AssertExpectations(t)
		})
	}
}