package repository

import (
	"context"
	"jobsite/models"
	"jobsite/request"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateJob(t *testing.T) {

	tests := []struct {
		name        string
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "jobs"`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "db error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "jobs"`).
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

			repo := NewJobRepository(db)

			job := &models.Job{
				Base: models.Base{
					ID: uuid.New(),
				},
				JobTitle:    "Backend Developer",
				CompanyName: "TestCorp",
			}

			err := repo.CreateJob(context.Background(), job)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetJob(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
		expectCount int64
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				jobID := uuid.New()
				jobTypeID := uuid.New()
				locationID := uuid.New()
				countryID := uuid.New()

				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				idRows := sqlmock.NewRows([]string{"id"}).AddRow(jobID)

				jobRows := sqlmock.NewRows([]string{
					"id", "website", "company_name", "job_title", "email", "description", "created_at",
				}).AddRow(jobID, "", "TestCorp", "Backend Developer", "", "", time.Now())

				mock.ExpectQuery(`SELECT COUNT\(DISTINCT`).WillReturnRows(countRows)
				mock.ExpectQuery(`SELECT jobs.id`).WillReturnRows(idRows)
				mock.ExpectQuery(`SELECT .* FROM "jobs"`).WillReturnRows(jobRows)

				mock.ExpectQuery(`SELECT .* FROM "job_job_types"`).
					WillReturnRows(
						sqlmock.NewRows([]string{"job_id", "job_type_id"}).
							AddRow(jobID, jobTypeID),
					)

				mock.ExpectQuery(`SELECT .* FROM "job_types"`).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "type_name"}).
							AddRow(jobTypeID, "Full Time"),
					)

				mock.ExpectQuery(`SELECT .* FROM "job_locations"`).
					WillReturnRows(
						sqlmock.NewRows([]string{"job_id", "location_id"}).
							AddRow(jobID, locationID),
					)

				mock.ExpectQuery(`SELECT .* FROM "locations"`).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "location_name", "country_id"}).
							AddRow(locationID, "Chennai", countryID),
					)

				mock.ExpectQuery(`SELECT .* FROM "countries"`).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "country_name"}).
							AddRow(countryID, "India"),
					)
			},
			expectCount: 1,
		},
		{
			name: "count error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT`).
					WillReturnError(assert.AnError)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			db, mock := SetupMockDB(t)

			tt.mockSetup(mock)

			repo := NewJobRepository(db)

			filter := request.JobFilter{

				Page: 1,

				PageSize: 10,
			}

			jobs, total, err := repo.GetJob(context.Background(), filter)

			if tt.expectError {

				assert.Error(t, err)

			} else {

				assert.NoError(t, err)

				assert.Equal(t, tt.expectCount, total)

				assert.NotNil(t, jobs)

			}
			assert.NoError(t, mock.ExpectationsWereMet())

		})

	}
}

func TestUpdateJob(t *testing.T) {

	tests := []struct {
		name        string
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE "jobs"`).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(`DELETE FROM "job_locations"`).WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectExec(`DELETE FROM "job_job_types"`).WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectCommit()
			},
		},
		{
			name: "update error",
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE "jobs"`).
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

			repo := NewJobRepository(db)

			job := &models.Job{
				Base: models.Base{
					ID: uuid.New(),
				},
				JobTitle:    "Updated Title",
				CompanyName: "Updated Company",
			}

			err := repo.UpdateJob(context.Background(), job)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteJob(t *testing.T) {

	tests := []struct {
		name        string
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {

				jobID := uuid.New()

				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(jobID)

				mock.ExpectQuery(`SELECT .* FROM "jobs"`).
					WillReturnRows(rows)

				mock.ExpectExec(`DELETE FROM "job_locations"`).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectExec(`DELETE FROM "job_job_types"`).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectExec(`DELETE FROM "jobs"`).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "job not found",
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT .* FROM "jobs"`).
					WillReturnError(gorm.ErrRecordNotFound)

				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			db, mock := SetupMockDB(t)
			tt.mockSetup(mock)

			repo := NewJobRepository(db)

			err := repo.DeleteJob(context.Background(), uuid.New())

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
