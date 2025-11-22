package entity

import (
	"time"
)

type LeaveRequestStatus string

const (
	Draft           LeaveRequestStatus = "draft"
	WaitingApproval LeaveRequestStatus = "waiting_approval"
	Approved        LeaveRequestStatus = "approved"
	Rejected        LeaveRequestStatus = "rejected"
)

func (r LeaveRequestStatus) IsValidStatus() bool {
	switch r {
	case Draft, WaitingApproval, Approved, Rejected:
		return true
	}
	return false
}

type LeaveRequestType string

const (
	Annual LeaveRequestType = "annual"
	Sick   LeaveRequestType = "sick"
	Unpaid LeaveRequestType = "unpaid"
)

func (r LeaveRequestType) IsValidType() bool {
	switch r {
	case Annual, Sick, Unpaid:
		return true
	}
	return false
}

type LeaveRequest struct {
	ID        int                `json:"id" db:"id"`
	UserId    int                `json:"userId" db:"user_id"`
	StartDate time.Time          `json:"startDate" db:"start_date"`
	EndDate   time.Time          `json:"endDate" db:"end_date"`
	Reason    string             `json:"reason" db:"reason"`
	Type      LeaveRequestType   `json:"type" db:"type"`
	Status    LeaveRequestStatus `json:"status" db:"status"`
	CreatedAt time.Time          `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt" db:"updated_at"`
}

type LeaveRequestFilter struct {
	UserId string
	Status string
}
