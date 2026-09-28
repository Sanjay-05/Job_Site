package handler

import (
	"bytes"
	"context"
	"errors"
	"jobsite/models"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type MockLookupService struct {
	Countries []models.Country
	JobTypes  []models.JobType
	Err       error
}

func (m *MockLookupService) ListCountries(ctx context.Context) ([]models.Country, error) {
	return m.Countries, m.Err
}

func (m *MockLookupService) ListJobTypes(ctx context.Context) ([]models.JobType, error) {
	return m.JobTypes, m.Err
}

func TestListCountriesEndpoint(t *testing.T) {

	tests := []struct {
		name      string
		mockSvc   *MockLookupService
		expectErr bool
	}{{
		name: "success",
		mockSvc: &MockLookupService{
			Countries: []models.Country{
				{Base: models.Base{ID: [16]byte{}}, CountryName: "USA"},
				{Base: models.Base{ID: [16]byte{}}, CountryName: "India"},
			},
		},
		expectErr: false,
	},
		{
			name: "service error",
			mockSvc: &MockLookupService{
				Err: errors.New("service error"),
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			e := &EndpointSet{
				LookupSvc: tt.mockSvc,
			}

			endpoint := e.ListCountries()

			resp, err := endpoint(context.Background(), nil)

			assert.NoError(t, err)

			if tt.expectErr {
				result, ok := resp.(GenericResponse)
				if !ok {
					t.Fatalf("expected GenericResponse, got %T", resp)
				}

				if result.Err == nil || result.Err.Error() != "service error" {
					t.Errorf("expected error 'service error', got %v", result.Err)
				}

			} else {
				result, ok := resp.(map[string]any)
				if !ok {
					t.Fatalf("expected map[string]any, got %T", resp)
				}

				data := result["data"].([]models.Country)

				if !reflect.DeepEqual(data, tt.mockSvc.Countries) {
					t.Errorf("expected %v, got %v", tt.mockSvc.Countries, data)
				}
			}
		})
	}
}

func TestListJobTypesEndpoint(t *testing.T) {

	tests := []struct {
		name      string
		mockSvc   *MockLookupService
		expectErr bool
	}{
		{
			name: "success",
			mockSvc: &MockLookupService{
				JobTypes: []models.JobType{
					{Base: models.Base{ID: [16]byte{}}, TypeName: "Full-time"},
					{Base: models.Base{ID: [16]byte{}}, TypeName: "Part-time"},
				},
			},
			expectErr: false,
		},
		{
			name: "service error",
			mockSvc: &MockLookupService{
				Err: errors.New("service error"),
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			e := &EndpointSet{
				LookupSvc: tt.mockSvc,
			}

			endpoint := e.ListJobTypes()

			resp, err := endpoint(context.Background(), nil)

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if tt.expectErr {

				result, ok := resp.(GenericResponse)
				if !ok {
					t.Fatalf("expected GenericResponse, got %T", resp)
				}

				if result.Err == nil || result.Err.Error() != "service error" {
					t.Errorf("expected error 'service error', got %v", result.Err)
				}

			} else {

				result, ok := resp.(map[string]any)
				if !ok {
					t.Fatalf("expected map[string]any, got %T", resp)
				}
				data := result["data"].([]models.JobType)
				if !reflect.DeepEqual(data, tt.mockSvc.JobTypes) {
					t.Errorf("expected %v, got %v", tt.mockSvc.JobTypes, data)
				}
			}
		})
	}
}

func TestDecodeRequests(t *testing.T) {

	tests := []struct {
		name       string
		method     string
		url        string
		body       string
		setVars    map[string]string
		decodeFunc func(context.Context, *http.Request) (any, error)
		expectErr  bool
	}{
		{
			name:       "decode register success",
			method:     "POST",
			url:        "/",
			body:       `{"username":"test","password":"123"}`,
			decodeFunc: DecodeRegisterUserRequest,
		},
		{
			name:       "decode login success",
			method:     "POST",
			url:        "/",
			body:       `{"username":"test","password":"123"}`,
			decodeFunc: DecodeLoginRequest,
		},
		{
			name:       "decode forget password success",
			method:     "POST",
			url:        "/",
			body:       `{"username":"test"}`,
			decodeFunc: DecodeForgetPasswordRequest,
		},
		{
			name:       "decode reset password success",
			method:     "POST",
			url:        "/",
			body:       `{"username":"test","otp":"123456","password":"newpass"}`,
			decodeFunc: DecodeResetPasswordRequest,
		},
		{
			name:       "decode post job success",
			method:     "POST",
			url:        "/",
			body:       `{"title":"golang dev","company_name":"abc"}`,
			decodeFunc: DecodePostJobRequest,
		},
		{
			name:       "decode update job success",
			method:     "PUT",
			url:        "/",
			body:       `{"title":"updated job"}`,
			decodeFunc: DecodeUpdateJobRequest,
		},
		{
			name:       "decode list jobs with query",
			method:     "GET",
			url:        "/?page=2&page_size=20&country=India",
			decodeFunc: DecodeListJobsRequest,
		},
		{
			name:       "decode empty request",
			method:     "GET",
			url:        "/",
			decodeFunc: DecodeEmptyRequest,
		},
		{
			name:   "decode delete job success",
			method: "DELETE",
			url:    "/job/1",
			setVars: map[string]string{
				"id": uuid.New().String(),
			},
			decodeFunc: DecodeDeleteJobRequest,
		},
		{
			name:       "decode delete job invalid uuid",
			method:     "DELETE",
			url:        "/job/1",
			setVars:    map[string]string{"id": "invalid"},
			decodeFunc: DecodeDeleteJobRequest,
			expectErr:  true,
		},
		{
			name:       "decode invalid json",
			method:     "POST",
			url:        "/",
			body:       `{invalid}`,
			decodeFunc: DecodeLoginRequest,
			expectErr:  true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(tt.method, tt.url, bytes.NewBufferString(tt.body))

			if tt.setVars != nil {
				req = mux.SetURLVars(req, tt.setVars)
			}

			_, err := tt.decodeFunc(context.Background(), req)

			if tt.expectErr && err == nil {
				t.Fatalf("expected error but got nil")
			}

			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error %v", err)
			}
		})
	}
}

func TestEncodeResponse(t *testing.T) {

	tests := []struct {
		name       string
		response   any
		expectCode int
	}{
		{
			name:       "success response",
			response:   AuthResponse{Message: "ok"},
			expectCode: http.StatusOK,
		},
		{
			name: "error response",
			response: GenericResponse{
				Err: errors.New("failed"),
			},
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			EncodeResponse(context.Background(), w, tt.response)

			if w.Code != tt.expectCode {
				t.Fatalf("expected %d got %d", tt.expectCode, w.Code)
			}
		})
	}
}
