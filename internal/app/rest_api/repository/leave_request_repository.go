package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/devonLoen/leave-request-service/internal/app/rest_api/database"
	entity "github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"
)

type LeaveRequest struct {
	database.BaseSQLRepository[entity.LeaveRequest]
}

func NewLeaveRequestRepository(db *sql.DB) *LeaveRequest {
	return &LeaveRequest{
		BaseSQLRepository: database.BaseSQLRepository[entity.LeaveRequest]{DB: db},
	}
}

func mapLeaveRequest(rows *sql.Row, lr *entity.LeaveRequest) error {
	return rows.Scan(&lr.ID, &lr.UserId, &lr.StartDate, &lr.EndDate, &lr.Type, &lr.Status, &lr.Reason)
}

func mapLeaveRequests(rows *sql.Rows, lr *entity.LeaveRequest) error {
	return rows.Scan(&lr.ID, &lr.StartDate, &lr.EndDate, &lr.Type, &lr.Status, &lr.Reason)
}

func (r *LeaveRequest) FindById(id int) (*entity.LeaveRequest, error) {
	return r.SelectSingle(
		mapLeaveRequest,
		"SELECT lr.id, lr.user_id, lr.start_date, lr.end_date, lr.type, lr.status, lr.reason FROM leave_requests lr WHERE lr.id = $1",
		id,
	)
}

func (r *LeaveRequest) GetAllLeaveRequests(limit, offset int, sortBy, orderBy, search string, filter entity.LeaveRequestFilter) ([]*entity.LeaveRequest, error) {
	baseQuery := "SELECT lr.id, lr.start_date, lr.end_date, lr.type, lr.status, lr.reason FROM leave_requests lr"
	var conditions []string
	var args []interface{}

	argId := 1

	if filter.UserId != "" {
		conditions = append(conditions, fmt.Sprintf("(lr.user_id = $%d)", argId))
		args = append(args, filter.UserId)
		argId++
	}

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("(lr.status = $%d)", argId))
		args = append(args, filter.Status)
		argId++
	}

	if search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(lr.type::text ILIKE $%d OR lr.status::text ILIKE $%d OR lr.reason ILIKE $%d)",
			argId, argId, argId,
		))
		args = append(args, "%"+search+"%")
		argId++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += fmt.Sprintf(" ORDER BY lr.%s %s LIMIT $%d OFFSET $%d", sortBy, orderBy, argId, argId+1)
	fmt.Println((baseQuery))
	args = append(args, limit, offset)

	return r.SelectMultiple(
		mapLeaveRequests,
		baseQuery,
		args...,
	)
}

func (r *LeaveRequest) Create(leaveRequest *entity.LeaveRequest) error {
	_, err := r.Insert(
		"INSERT INTO leave_requests (user_id, start_date, end_date, type, status, reason) VALUES ($1, $2, $3, $4, $5, $6)",
		leaveRequest.UserId, leaveRequest.StartDate, leaveRequest.EndDate, leaveRequest.Type, leaveRequest.Status, leaveRequest.Reason,
	)
	return err
}

func (r *LeaveRequest) Approve(leaveRequestId int) error {
	_, err := r.ExecuteQuery(
		"UPDATE leave_requests SET status = 'approved' WHERE id = $1",
		leaveRequestId,
	)
	return err
}

func (r *LeaveRequest) Reject(leaveRequestId int) error {
	_, err := r.ExecuteQuery(
		"UPDATE leave_requests SET status = 'rejected' WHERE id = $1",
		leaveRequestId,
	)
	return err
}

func (r *LeaveRequest) OverlapApprovedLeaveExists(userId int, startDate, endDate time.Time) (bool, error) {
	var exists bool

	query := `
        SELECT EXISTS (
            SELECT 1 
            FROM leave_requests lr
            WHERE lr.user_id = $1 
            AND lr.status = 'approved' 
            AND NOT (lr.end_date < $2 OR lr.start_date > $3)
        )`

	row := r.DB.QueryRow(query, userId, startDate, endDate)

	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
