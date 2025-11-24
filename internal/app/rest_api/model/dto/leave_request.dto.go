package dto

import (
	"time"

	entity "github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"
)

type LeaveRequestResponse struct {
	ID        int       `json:"id"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Reason    string    `json:"reason"`
}

type GetAllLeaveRequestsResponse struct {
	LeaveRequests []*LeaveRequestResponse `json:"leaveRequests"`
}

type CreateLeaveRequestRequest struct {
	StartDate time.Time `json:"startDate" validate:"required"`
	EndDate   time.Time `json:"endDate" validate:"required"`
	Type      string    `json:"type" validate:"required,oneof=annual sick unpaid"`
	Reason    string    `json:"reason" validate:"required,min=10,max=500"`
	Status    string    `json:"status" validate:"required,oneof=draft waiting_approval"`
}

type CreateLeaveRequestResponse struct {
	StartDate time.Time `json:"startDate" validate:"required"`
	EndDate   time.Time `json:"endDate" validate:"required"`
	Type      string    `json:"type" validate:"required,oneof=annual sick unpaid"`
	Reason    string    `json:"reason" validate:"required,min=10,max=500"`
	Status    string    `json:"status" validate:"required,oneof=draft waiting_approval"`
	Message   string    `json:"message" binding:"required"`
}

func (r *GetAllLeaveRequestsResponse) MapLeaveRequestsResponse(leaveRequests []*entity.LeaveRequest) {
	for _, leaveRequests := range leaveRequests {
		leaveRequest := &LeaveRequestResponse{
			ID:        leaveRequests.ID,
			StartDate: leaveRequests.StartDate,
			EndDate:   leaveRequests.EndDate,
			Type:      string(leaveRequests.Type),
			Status:    string(leaveRequests.Status),
			Reason:    leaveRequests.Reason,
		}
		r.LeaveRequests = append(r.LeaveRequests, leaveRequest)
	}
}

func (r *LeaveRequestResponse) MapLeaveRequestResponse(leaveRequest *entity.LeaveRequest) {
	r.ID = leaveRequest.ID
	r.StartDate = leaveRequest.StartDate
	r.EndDate = leaveRequest.EndDate
	r.Type = string(leaveRequest.Type)
	r.Status = string(leaveRequest.Status)
	r.Reason = leaveRequest.Reason
}

func (ur *CreateLeaveRequestRequest) ToLeaveRequest(userId int) *entity.LeaveRequest {
	return &entity.LeaveRequest{
		UserId:    userId,
		StartDate: ur.StartDate,
		EndDate:   ur.EndDate,
		Type:      entity.LeaveRequestType(ur.Type),
		Status:    entity.LeaveRequestStatus(ur.Status),
		Reason:    ur.Reason,
	}
}

func (ur *CreateLeaveRequestResponse) FromLeaveRequest(leaveRequest *entity.LeaveRequest) *CreateLeaveRequestResponse {
	return &CreateLeaveRequestResponse{
		StartDate: leaveRequest.StartDate,
		EndDate:   leaveRequest.EndDate,
		Type:      string(leaveRequest.Type),
		Status:    string(leaveRequest.Status),
		Reason:    leaveRequest.Reason,
		Message:   "Leave Request created successfully.",
	}
}
