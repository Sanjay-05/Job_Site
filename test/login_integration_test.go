package test

import (
	"context"
	"os"
	"testing"

	"jobsite/handler"
	"jobsite/models"
	"jobsite/repository"
	"jobsite/request"
	"jobsite/service"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestLoginIntegration(t *testing.T) {

	ctx := context.Background()
	db := setupTestDB(t)

	userRepo := repository.NewUserRepository(db)
	logger := log.NewLogfmtLogger(os.Stdout)

	authSvc := service.NewUserService(userRepo, logger, "test-secret")
	endpoints := handler.NewEndpoints(authSvc, nil, nil)

	loginEndpoint := endpoints.Login()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := models.User{
		UserName: "test_user",
		Password: string(hashedPassword),
		Role:     "user",
		Email:    "test@example.com",
		PhoneNo:  "9876543210",
	}

	err = db.Create(&user).Error
	require.NoError(t, err)

	tests := []struct {
		name        string
		req         request.LoginInput
		expectErr   bool
		expectMsg   string
		expectToken bool
	}{
		{
			name: "valid login",
			req: request.LoginInput{
				UserName: "test_user",
				Password: "user123",
			},
			expectErr:   false,
			expectMsg:   "Login successful",
			expectToken: true,
		},
		{
			name: "invalid password",
			req: request.LoginInput{
				UserName: "test_user",
				Password: "wrong",
			},
			expectErr: true,
			expectMsg: "invalid credentials",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			resp, err := loginEndpoint(ctx, tt.req)
			require.NoError(t, err)

			if tt.expectErr {
				r := resp.(handler.GenericResponse)
				assert.Error(t, r.Err)
				assert.Equal(t, tt.expectMsg, r.Err.Error())
			} else {
				r := resp.(handler.AuthResponse)
				assert.Equal(t, tt.expectMsg, r.Message)

				if tt.expectToken {
					assert.NotEmpty(t, r.Token)
				}
			}
		})
	}
}

func setupTestDB(t *testing.T) *gorm.DB {

	dsn := "host=localhost user=postgres password=password dbname=jobDB port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	db.Exec("TRUNCATE TABLE users RESTART IDENTITY")

	return db
}
