package handler_test

import (
	"context"
	"errors"
	"jobsite/handler"
	"jobsite/request"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(ctx context.Context, req request.RegisterUserInput) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockUserService) Login(ctx context.Context, userName, password string) (string, error) {
	args := m.Called(ctx, userName, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) ForgetPassword(ctx context.Context, userName string) error {
	args := m.Called(ctx, userName)
	return args.Error(0)
}

func (m *MockUserService) ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func TestLoginEndpoint(t *testing.T) {

	tests := []struct {
		name        string
		username    string
		password    string
		token       string
		serviceErr  error
		expectError bool
	}{
		{name: "Success", username: "test_user", password: "user123", token: "jwt-token", serviceErr: nil, expectError: false},
		{name: "invalid credentials", username: "john", password: "wrong", token: "", serviceErr: errors.New("invalid"), expectError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockUserService)
			mockSvc.On("Login", mock.Anything, tt.username, tt.password).
				Return(tt.token, tt.serviceErr)
			endpoints := handler.NewEndpoints(mockSvc, nil, nil)
			ep := endpoints.Login()

			response, err := ep(context.Background(), request.LoginInput{
				UserName: tt.username,
				Password: tt.password,
			})
			assert.NoError(t, err)
			if tt.expectError {
				r := response.(handler.GenericResponse)
				assert.Error(t, r.Err)
			} else {
				r := response.(handler.AuthResponse)
				assert.Equal(t, "Login successful", r.Message)
				assert.Equal(t, tt.token, r.Token)
			}
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestRegisterUserEndpoint(t *testing.T) {
	tests := []struct {
		name        string
		req         request.RegisterUserInput
		serviceErr  error
		expectError bool
		expectedMsg string
	}{
		{
			name: "success",
			req: request.RegisterUserInput{
				UserName: "test_user1",
				Password: "User@123",
			},
			serviceErr:  nil,
			expectError: false,
			expectedMsg: "Registration Successful",
		},
		{
			name: "service error",
			req: request.RegisterUserInput{
				UserName: "test_user1",
				Password: "User@123",
			},
			serviceErr:  errors.New("registration failed"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockUserService)
			mockSvc.On("RegisterUser", mock.Anything, tt.req).Return(tt.serviceErr)

			endpoints := handler.NewEndpoints(mockSvc, nil, nil)
			ep := endpoints.RegisterUser()

			resp, err := ep(context.Background(), tt.req)
			assert.NoError(t, err)

			if tt.expectError {
				r := resp.(handler.GenericResponse)
				assert.Error(t, r.Err)
			} else {
				r := resp.(handler.AuthResponse)
				assert.Equal(t, tt.expectedMsg, r.Message)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestForgetPasswordEndpoint(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		serviceErr  error
		expectError bool
	}{
		{
			name:        "success",
			username:    "john",
			serviceErr:  nil,
			expectError: false,
		},
		{
			name:        "service error",
			username:    "john",
			serviceErr:  errors.New("user not found"),
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockUserService)
			mockSvc.On("ForgetPassword", mock.Anything, tt.username).
				Return(tt.serviceErr)
			endpoints := handler.NewEndpoints(mockSvc, nil, nil)
			ep := endpoints.ForgetPassword()
			response, err := ep(context.Background(), request.ForgetPasswordRequest{
				UserName: tt.username,
			})
			assert.NoError(t, err)
			if tt.expectError {
				r := response.(handler.GenericResponse)
				assert.Error(t, r.Err)
			} else {
				r := response.(handler.AuthResponse)
				assert.Equal(t, "Otp send Successfull", r.Message)
			}
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestResetPasswordEndpoint(t *testing.T) {
	tests := []struct {
		name        string
		req         request.ResetPasswordRequest
		serviceErr  error
		expectError bool
	}{
		{
			name: "success",
			req: request.ResetPasswordRequest{
				Username:           "john",
				OTP:                "123456",
				NewPassword:        "New@123",
				ConfirmNewPassword: "New@123",
			},
			serviceErr:  nil,
			expectError: false,
		},
		{
			name: "invalid otp",
			req: request.ResetPasswordRequest{
				Username:           "john",
				OTP:                "000000",
				NewPassword:        "New@123",
				ConfirmNewPassword: "New@123",
			},
			serviceErr:  errors.New("invalid otp"),
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockUserService)
			mockSvc.On("ResetPassword", mock.Anything, &tt.req).
				Return(tt.serviceErr)
			endpoints := handler.NewEndpoints(mockSvc, nil, nil)
			ep := endpoints.ResetPassword()
			response, err := ep(context.Background(), tt.req)
			assert.NoError(t, err)
			if tt.expectError {
				r := response.(handler.GenericResponse)
				assert.Error(t, r.Err)
			} else {
				r := response.(handler.AuthResponse)
				assert.Equal(t, "password reset successfull", r.Message)
			}
			mockSvc.AssertExpectations(t)
		})
	}
}
