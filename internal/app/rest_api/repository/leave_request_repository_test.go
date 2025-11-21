package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var ErrSimulatedDB = errors.New("simulated DB error")

func setupMockDB(t *testing.T) (*LeaveRequest, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}

	repo := &LeaveRequest{DB: db}
	return repo, mock
}

func TestOverlapApprovedLeaveExists(t *testing.T) {
	userID := 101

	startDate := time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, time.February, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		expectedExist bool
		mockExpect    func(mock sqlmock.Sqlmock, userID int, start, end time.Time)
		expectedError error
	}{
		{
			name:          "Case 1: Overlap Found",
			expectedExist: true,
			expectedError: nil,
			mockExpect: func(mock sqlmock.Sqlmock, userID int, start, end time.Time) {
				mock.ExpectQuery(`SELECT EXISTS`).
					WithArgs(userID, start, end).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
		},
		{
			name:          "Case 2: No Overlap Found",
			expectedExist: false,
			expectedError: nil,
			mockExpect: func(mock sqlmock.Sqlmock, userID int, start, end time.Time) {
				mock.ExpectQuery(`SELECT EXISTS`).
					WithArgs(userID, start, end).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
		},
		{
			name:          "Case 3: Database Error Occurred",
			expectedExist: false,
			expectedError: ErrSimulatedDB,
			mockExpect: func(mock sqlmock.Sqlmock, userID int, start, end time.Time) {
				mock.ExpectQuery(`SELECT EXISTS`).
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
