package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"
	handler "github.com/devonLoen/leave-request-service/internal/app/rest_api/handler"
	models "github.com/devonLoen/leave-request-service/internal/app/rest_api/model"
	dto "github.com/devonLoen/leave-request-service/internal/app/rest_api/model/dto"
)

type MockLeaveRequestUsecase struct {
	mock.Mock
}

func (m *MockLeaveRequestUsecase) CreateLeaveRequest(req *dto.CreateLeaveRequestRequest, userID int) (*dto.CreateLeaveRequestResponse, *models.ErrorResponse) {
	args := m.Called(req, userID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*models.ErrorResponse)
	}
	return args.Get(0).(*dto.CreateLeaveRequestResponse), nil
}

func (m *MockLeaveRequestUsecase) GetAllLeaveRequests(
	limit, offset int,
	sortBy, orderBy, search string,
	filter entity.LeaveRequestFilter,
) (*dto.GetAllLeaveRequestsResponse, *models.ErrorResponse) {
	args := m.Called(limit, offset, sortBy, orderBy, search, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*models.ErrorResponse)
	}
	return args.Get(0).(*dto.GetAllLeaveRequestsResponse), nil
}

func (m *MockLeaveRequestUsecase) GetLeaveRequest(id int) (*dto.LeaveRequestResponse, *models.ErrorResponse) {
	return nil, nil
}

func (m *MockLeaveRequestUsecase) Approve(id int) *models.ErrorResponse        { return nil }
func (m *MockLeaveRequestUsecase) Reject(id int) *models.ErrorResponse         { return nil }
func (m *MockLeaveRequestUsecase) Submit(id, userID int) *models.ErrorResponse { return nil }

func TestCreateLeaveRequestHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	today := time.Now().Truncate(24 * time.Hour)

	tests := []struct {
		name           string
		input          dto.CreateLeaveRequestRequest
		mockSetup      func(m *MockLeaveRequestUsecase)
		expectedCode   int
		expectedErrMsg string
	}{
		{
			name: "Start date after end date",
			input: dto.CreateLeaveRequestRequest{
				StartDate: today.Add(48 * time.Hour),
				EndDate:   today.Add(24 * time.Hour),
				Reason:    "Annual leave",
				Type:      "annual",
				Status:    "waiting_approval",
			},
			mockSetup:      func(m *MockLeaveRequestUsecase) {},
			expectedCode:   http.StatusBadRequest,
			expectedErrMsg: "startDate cannot be after endDate",
		},
		{
			name: "Start date in the past",
			input: dto.CreateLeaveRequestRequest{
				StartDate: today.Add(-24 * time.Hour),
				EndDate:   today.Add(24 * time.Hour),
				Reason:    "Annual leave",
				Type:      "annual",
				Status:    "waiting_approval",
			},
			mockSetup:      func(m *MockLeaveRequestUsecase) {},
			expectedCode:   http.StatusBadRequest,
			expectedErrMsg: "Leave request cannot be in the past",
		},
		{
			name: "Usecase error",
			input: dto.CreateLeaveRequestRequest{
				StartDate: today.Add(24 * time.Hour),
				EndDate:   today.Add(48 * time.Hour),
				Reason:    "Annual leave",
				Type:      "annual",
				Status:    "waiting_approval",
			},
			mockSetup: func(m *MockLeaveRequestUsecase) {
				m.On("CreateLeaveRequest", mock.Anything, 1).
					Return(nil, &models.ErrorResponse{
						Code:    http.StatusInternalServerError,
						Message: "Failed to create leave request",
					}).Once()
			},
			expectedCode:   http.StatusInternalServerError,
			expectedErrMsg: "Failed to create leave request",
		},
		{
			name: "Success",
			input: dto.CreateLeaveRequestRequest{
				StartDate: today.Add(24 * time.Hour),
				EndDate:   today.Add(48 * time.Hour),
				Reason:    "Annual leave",
				Type:      "annual",
				Status:    "waiting_approval",
			},
			mockSetup: func(m *MockLeaveRequestUsecase) {
				m.On("CreateLeaveRequest", mock.Anything, 1).
					Return(&dto.CreateLeaveRequestResponse{}, nil).Once()
			},
			expectedCode:   http.StatusCreated,
			expectedErrMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := new(MockLeaveRequestUsecase)
			tt.mockSetup(mockUC)

			r := gin.New()
			h := handler.NewLeaveRequestHandler(mockUC)
			r.POST("/leave", func(c *gin.Context) {
				c.Set("userId", 1)
				h.CreateLeaveRequest(c)
			})

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest(http.MethodPost, "/leave", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedErrMsg != "" {
				assert.Contains(t, w.Body.String(), tt.expectedErrMsg)
			}

			mockUC.AssertExpectations(t)
		})
	}
}
