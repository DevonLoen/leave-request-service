package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	entity "github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"
	models "github.com/devonLoen/leave-request-service/internal/app/rest_api/model"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/model/dto"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/repository"
)

type LeaveRequestUsecase interface {
	CreateLeaveRequest(*dto.CreateLeaveRequestRequest, int) (*dto.CreateLeaveRequestResponse, *models.ErrorResponse)
	GetAllLeaveRequests(limit, offset int, sortBy, orderBy, search string, filter entity.LeaveRequestFilter) (*dto.GetAllLeaveRequestsResponse, *models.ErrorResponse)
	GetLeaveRequest(leaveRequestID int) (*dto.LeaveRequestResponse, *models.ErrorResponse)
	Approve(leaveRequestID int) *models.ErrorResponse
	Reject(leaveRequestID int) *models.ErrorResponse
	Submit(leaveRequestID, userID int) *models.ErrorResponse
}

type LeaveRequest struct {
	leaveRequestRepo repository.LeaveRequestRepository
}

func NewLeaveRequestUsecase(leaveRequestRepo repository.LeaveRequestRepository) *LeaveRequest {
	return &LeaveRequest{leaveRequestRepo: leaveRequestRepo}
}

func (us *LeaveRequest) GetAllLeaveRequests(limit, offset int, sortBy, orderBy, search string, filter entity.LeaveRequestFilter) (*dto.GetAllLeaveRequestsResponse, *models.ErrorResponse) {
	response := &dto.GetAllLeaveRequestsResponse{}

	allowedSorts := map[string]bool{
		"id":        true,
		"userId":    true,
		"startDate": true,
		"endDate":   true,
		"type":      true,
		"status":    true,
		"reason":    true,
	}

	safeSortBy := "id"
	if allowedSorts[sortBy] {
		safeSortBy = sortBy
	} else {
		return nil, &models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid sort parameter",
		}
	}

	safeOrderBy := "ASC"
	if strings.ToUpper(orderBy) == "DESC" {
		safeOrderBy = "DESC"
	} else if strings.ToUpper(orderBy) != "ASC" {
		return nil, &models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid order parameter",
		}
	}

	if filter.Status != "" {
		statusEnum := entity.LeaveRequestStatus(filter.Status)
		if !statusEnum.IsValidStatus() {
			return nil, &models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid Status Filter parameter",
			}
		}
	}
	fmt.Println(("here?"))
	queriedLeaveRequests, err := us.leaveRequestRepo.GetAllLeaveRequests(limit, offset, safeSortBy, safeOrderBy, search, filter)
	if err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	response.MapLeaveRequestsResponse(queriedLeaveRequests)

	return response, nil
}

func (us *LeaveRequest) GetLeaveRequest(leaveRequestID int) (*dto.LeaveRequestResponse, *models.ErrorResponse) {
	response := &dto.LeaveRequestResponse{}

	leaveRequest, err := us.leaveRequestRepo.FindById(leaveRequestID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Leave Request Not Found",
			}
		}
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	response.MapLeaveRequestResponse(leaveRequest)

	return response, nil
}

func (us *LeaveRequest) CreateLeaveRequest(createLeaveRequestRequest *dto.CreateLeaveRequestRequest, userId int) (*dto.CreateLeaveRequestResponse, *models.ErrorResponse) {
	leaveRequestResponse := &dto.CreateLeaveRequestResponse{}
	errCheckExist := us.OverlapApprovedLeaveExists(userId, createLeaveRequestRequest.StartDate, createLeaveRequestRequest.EndDate)
	if errCheckExist != nil {
		return nil, errCheckExist
	}

	leaveRequest := createLeaveRequestRequest.ToLeaveRequest(userId)

	err := us.leaveRequestRepo.Create(leaveRequest)
	if err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create leave Request",
		}
	}

	return leaveRequestResponse.FromLeaveRequest(leaveRequest), nil
}

func (us *LeaveRequest) Approve(leaveRequestID int) *models.ErrorResponse {
	existingLeaveRequest, err := us.leaveRequestRepo.FindById(leaveRequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Leave Request not found",
			}
		}
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	if existingLeaveRequest.Status == "draft" {
		return &models.ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: "Leave Request status is still draft",
		}
	}

	if existingLeaveRequest.Status == "approved" {
		return &models.ErrorResponse{
			Code:    http.StatusConflict,
			Message: "Leave Request has been approved",
		}
	}

	errCheckExist := us.OverlapApprovedLeaveExists(existingLeaveRequest.UserId, existingLeaveRequest.StartDate, existingLeaveRequest.EndDate)
	if errCheckExist != nil {
		return errCheckExist
	}

	err = us.leaveRequestRepo.Approve(existingLeaveRequest.ID)

	if err != nil {
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to Approve Leave Request",
		}
	}

	return nil
}

func (us *LeaveRequest) Reject(leaveRequestID int) *models.ErrorResponse {
	existingLeaveRequest, err := us.leaveRequestRepo.FindById(leaveRequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Leave Request not found",
			}
		}
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	if existingLeaveRequest.Status == "draft" {
		return &models.ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: "Leave Request status is still draft",
		}
	}

	err = us.leaveRequestRepo.Reject(existingLeaveRequest.ID)

	if err != nil {
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to Reject Leave Request",
		}
	}

	return nil
}

func (lr *LeaveRequest) OverlapApprovedLeaveExists(userId int, startDate, endDate time.Time) *models.ErrorResponse {
	isOverlapping, err := lr.leaveRequestRepo.OverlapApprovedLeaveExists(userId, startDate, endDate)
	if err != nil {
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	if isOverlapping {
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "The requested leave dates overlap with an already approved leave request.",
		}
	}
	return nil
}

func (us *LeaveRequest) Submit(leaveRequestID, userId int) *models.ErrorResponse {
	existingLeaveRequest, err := us.leaveRequestRepo.FindById(leaveRequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Leave Request not found",
			}
		}
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	if existingLeaveRequest.UserId != userId {
		return &models.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "The specified leave request belongs to another user.",
		}
	}

	if existingLeaveRequest.Status != "draft" {
		return &models.ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: "Only Draft leave Request can be submitted",
		}
	}

	err = us.leaveRequestRepo.Submit(existingLeaveRequest.ID)

	if err != nil {
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to Submit Leave Request",
		}
	}

	return nil
}
