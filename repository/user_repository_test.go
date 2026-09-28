package repository

import (
	"context"
	"jobsite/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	assert.NoError(t, err)

	return gormDB, mock
}

func TestFindUserByName(t *testing.T) {

	tests := []struct {
		name        string
		username    string
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name:     "success",
			username: "john",
			mockSetup: func(mock sqlmock.Sqlmock) {
				id := uuid.New()

				rows := sqlmock.NewRows([]string{"id", "user_name", "password"}).
					AddRow(id, "john", "hashedpass")

				mock.ExpectQuery(`SELECT .* FROM "users"`).
					WithArgs("john", 1).
					WillReturnRows(rows)
			},
			expectError: false,
		},
		{
			name:     "user not found",
			username: "john",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT .* FROM "users"`).
					WithArgs("john", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			db, mock := SetupMockDB(t)

			tt.mockSetup(mock)

			repo := NewUserRepository(db)

			user, err := repo.FindUserByName(context.Background(), tt.username)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.username, user.UserName)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateUser(t *testing.T) {

	tests := []struct {
		name        string
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`INSERT INTO "users"`).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "db error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`INSERT INTO "users"`).
					WillReturnError(assert.AnError)

				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			db, mock := SetupMockDB(t)

			tt.mockSetup(mock)

			repo := NewUserRepository(db)

			user := &models.User{
				UserName: "john",
				Password: "hashedpass",
			}

			err := repo.CreateUser(context.Background(), user)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestResetPassword(t *testing.T) {

	tests := []struct {
		name        string
		rows        int64
		expectError bool
	}{
		{
			name:        "success",
			rows:        1,
			expectError: false,
		},
		{
			name:        "no rows affected",
			rows:        0,
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			db, mock := SetupMockDB(t)

			mock.ExpectBegin()

			mock.ExpectExec(`UPDATE "users"`).
				WillReturnResult(sqlmock.NewResult(1, tt.rows))

			mock.ExpectCommit()

			repo := NewUserRepository(db)

			err := repo.ResetPassword(context.Background(), "john", "newpass")

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
