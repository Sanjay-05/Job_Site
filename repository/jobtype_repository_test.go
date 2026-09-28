package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJobTypeRepository_List(t *testing.T) {

	tests := []struct {
		name        string
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
		expectCount int
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {

				id1 := uuid.New()
				id2 := uuid.New()

				rows := sqlmock.NewRows([]string{"id", "type_name"}).
					AddRow(id1, "Full Time").
					AddRow(id2, "Part Time")

				mock.ExpectQuery(`SELECT .* FROM "job_types"`).
					WillReturnRows(rows)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "db error",
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(`SELECT .* FROM "job_types"`).
					WillReturnError(assert.AnError)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			db, mock := SetupMockDB(t)

			tt.mockSetup(mock)

			repo := NewJobTypeRepository(db)

			result, err := repo.List(context.Background())

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectCount)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
