package repository

import (
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/database"
	entity "github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"
)

var ErrSimulatedDB = errors.New("simulated DB error")

var leaveRequestColumns = []string{
	"id", "user_id", "start_date", "end_date", "type", "status", "reason",
}

func setupMockDB(t *testing.T) (*LeaveRequest, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}

	repo := &LeaveRequest{
		DB: db,

		BaseSQLRepository: database.BaseSQLRepository[entity.LeaveRequest]{
			DB: db,
		},
	}
	return repo, mock
}

func TestOverlapApprovedLeaveExists(t *testing.T) {
	userID := 101

	startDate := time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, time.February, 10, 0, 0, 0, 0, time.UTC)

	var overlappingRow = []driver.Value{
		1, userID,
		time.Date(2025, time.February, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2025, time.February, 15, 0, 0, 0, 0, time.UTC),
		"ANNUAL", "approved", "Holiday",
	}

	tests := []struct {
		name          string
		expectedExist bool
		mockExpect    func(mock sqlmock.Sqlmock, userID int, start, end time.Time)
		expectedError error
	}{
		{
			name:          "Case 1: Overlap Found (Returns 1 Row)",
			expectedExist: true,
			expectedError: nil,
			mockExpect: func(mock sqlmock.Sqlmock, userID int, start, end time.Time) {
				mock.ExpectQuery(`SELECT lr.id`).
					WithArgs(userID, start, end).
					WillReturnRows(
						sqlmock.NewRows(leaveRequestColumns).
							AddRow(overlappingRow...),
					)
			},
		},
		{
			name:          "Case 2: No Overlap Found (Returns 0 Rows)",
			expectedExist: false,
			expectedError: nil,
			mockExpect: func(mock sqlmock.Sqlmock, userID int, start, end time.Time) {
				mock.ExpectQuery(`SELECT lr.id`).
					WithArgs(userID, start, end).
					WillReturnRows(sqlmock.NewRows(leaveRequestColumns))
			},
		},
		{
			name:          "Case 3: Database Error Occurred",
			expectedExist: false,
			expectedError: ErrSimulatedDB,
			mockExpect: func(mock sqlmock.Sqlmock, userID int, start, end time.Time) {
				mock.ExpectQuery(`SELECT lr.id`).
					WithArgs(userID, start, end).
					WillReturnError(ErrSimulatedDB)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock := setupMockDB(t)
			defer repo.DB.Close()

			tt.mockExpect(mock, userID, startDate, endDate)

			exists, err := repo.OverlapApprovedLeaveExists(userID, startDate, endDate)

			if errMock := mock.ExpectationsWereMet(); errMock != nil {
				t.Fatalf("Mock expectations not met: %s", errMock)
			}

			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("Error mismatch. Actual: %v, Expected: %v", err, tt.expectedError)
			}

			if tt.expectedError != nil {
				return
			}

			if exists != tt.expectedExist {
				t.Errorf("Result 'exists' mismatch. Actual: %t, Expected: %t", exists, tt.expectedExist)
			}
		})
	}
}
