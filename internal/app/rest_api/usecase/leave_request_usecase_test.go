package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	entity "github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/model/dto"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/usecase"
)

type MockLeaveRequestRepo struct {
	mock.Mock
}

func (m *MockLeaveRequestRepo) OverlapApprovedLeaveExists(userId int, startDate, endDate time.Time) (bool, error) {
	args := m.Called(userId, startDate, endDate)
	return args.Bool(0), args.Error(1)
}

func (m *MockLeaveRequestRepo) Create(lr *entity.LeaveRequest) error {
	return m.Called(lr).Error(0)
}

func (m *MockLeaveRequestRepo) FindById(id int) (*entity.LeaveRequest, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.LeaveRequest), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLeaveRequestRepo) GetAllLeaveRequests(limit, offset int, sortBy, orderBy, search string, filter entity.LeaveRequestFilter) ([]*entity.LeaveRequest, error) {
	args := m.Called(limit, offset, sortBy, orderBy, search, filter)
	if args.Get(0) != nil {
		return args.Get(0).([]*entity.LeaveRequest), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLeaveRequestRepo) Approve(leaveRequestID int) error {
	return m.Called(leaveRequestID).Error(0)
}

func (m *MockLeaveRequestRepo) Reject(leaveRequestID int) error {
	return m.Called(leaveRequestID).Error(0)
}

func TestCreateLeaveRequest(t *testing.T) {

	mockRepo := new(MockLeaveRequestRepo)
	uc := usecase.NewLeaveRequestUsecase(mockRepo)

	today := time.Now().Truncate(24 * time.Hour)

	tests := []struct {
		name       string
		req        dto.CreateLeaveRequestRequest
		setupMock  func()
		wantErr    bool
		errMessage string
	}{
		{
			name: "Overlap detected",
			req: dto.CreateLeaveRequestRequest{
				StartDate: today.Add(24 * time.Hour),
				EndDate:   today.Add(48 * time.Hour),
			},
			setupMock: func() {
				mockRepo.On("OverlapApprovedLeaveExists", 1, mock.Anything, mock.Anything).
					Return(true, nil).Once()
			},
			wantErr:    true,
			errMessage: "The requested leave dates overlap with an already approved leave request.",
		},
		{
			name: "Repo create error",
			req: dto.CreateLeaveRequestRequest{
				StartDate: today.Add(24 * time.Hour),
				EndDate:   today.Add(48 * time.Hour),
			},
			setupMock: func() {
				mockRepo.On("OverlapApprovedLeaveExists", 1, mock.Anything, mock.Anything).
					Return(false, nil).Once()
				mockRepo.On("Create", mock.Anything).
					Return(errors.New("db error")).Once()
			},
			wantErr:    true,
			errMessage: "Failed to create leave Request",
		},
		{
			name: "Success create",
			req: dto.CreateLeaveRequestRequest{
				StartDate: today.Add(24 * time.Hour),
				EndDate:   today.Add(48 * time.Hour),
			},
			setupMock: func() {
				mockRepo.On("OverlapApprovedLeaveExists", 1, mock.Anything, mock.Anything).
					Return(false, nil).Once()
				mockRepo.On("Create", mock.Anything).
					Return(nil).Once()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.Mock.ExpectedCalls = nil

			tt.setupMock()

			res, errResp := uc.CreateLeaveRequest(&tt.req, 1)

			if tt.wantErr {
				assert.Nil(t, res)
				assert.NotNil(t, errResp)
				assert.Contains(t, errResp.Message, tt.errMessage, "Error message should contain expected text")
			} else {
				assert.NotNil(t, res)
				assert.Nil(t, errResp)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
