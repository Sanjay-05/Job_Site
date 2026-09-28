package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCountryRepository_List(t *testing.T) {

	tests := []struct {
		name        string
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
		expectLen   int
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {

				rows := sqlmock.NewRows([]string{
					"id",
					"country_name",
				}).
					AddRow(uuid.New(), "India").
					AddRow(uuid.New(), "USA")

				mock.ExpectQuery(`SELECT \* FROM "countries"`).
					WillReturnRows(rows)
			},
			expectLen: 2,
		},
		{
			name: "db error",
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(`SELECT \* FROM "countries"`).
					WillReturnError(assert.AnError)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			db, mock := SetupMockDB(t)
			tt.mockSetup(mock)

			repo := NewCountryRepository(db)

			countries, err := repo.List(context.Background())

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, countries, tt.expectLen)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
