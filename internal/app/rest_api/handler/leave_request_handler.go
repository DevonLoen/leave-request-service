package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"
	dto "github.com/devonLoen/leave-request-service/internal/app/rest_api/model/dto"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/pkg/util"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LeaveRequest struct {
	leaveRequestUsecase usecase.LeaveRequestUsecase
}

func NewLeaveRequestHandler(leaveRequestUsecase usecase.LeaveRequestUsecase) *LeaveRequest {
	return &LeaveRequest{leaveRequestUsecase: leaveRequestUsecase}
}

func (h *LeaveRequest) GetAllLeaveRequests(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	sortByStr := ctx.DefaultQuery("sortBy", "id")
	orderByStr := ctx.DefaultQuery("orderBy", "asc")

	filter := entity.LeaveRequestFilter{
		UserId: ctx.Query("userId"),
		Status: ctx.Query("status"),
	}

	search := ctx.Query("search")

	page, errConv := strconv.Atoi(pageStr)
	if errConv != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Page not valid "})

		return
	}

	limit, errConv := strconv.Atoi(limitStr)
	if errConv != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "limit not valid "})

		return
	}

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	allUsers, err := h.leaveRequestUsecase.GetAllLeaveRequests(limit, offset, sortByStr, orderByStr, search, filter)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code, err)

		return
	}

	ctx.JSON(http.StatusOK, allUsers)
}

func (h *LeaveRequest) GetMyLeaveRequests(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	sortByStr := ctx.DefaultQuery("sortBy", "id")
	orderByStr := ctx.DefaultQuery("orderBy", "asc")
	userIDRaw, _ := ctx.Get("userId")
	userID := userIDRaw.(int)

	filter := entity.LeaveRequestFilter{
		UserId: strconv.Itoa(userID),
		Status: ctx.Query("status"),
	}

	search := ctx.Query("search")

	page, errConv := strconv.Atoi(pageStr)
	if errConv != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Page not valid "})

		return
	}

	limit, errConv := strconv.Atoi(limitStr)
	if errConv != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "limit not valid "})

		return
	}

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	allUsers, err := h.leaveRequestUsecase.GetAllLeaveRequests(limit, offset, sortByStr, orderByStr, search, filter)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code, err)

		return
	}

	ctx.JSON(http.StatusOK, allUsers)
}

func (h *LeaveRequest) GetLeaveRequest(ctx *gin.Context) {
	leaveRequestID, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Leave Request ID not valid"})

		return
	}

	leaveRequest, leaveRequestErr := h.leaveRequestUsecase.GetLeaveRequest(leaveRequestID)
	if leaveRequestErr != nil {
		ctx.AbortWithStatusJSON(leaveRequestErr.Code, leaveRequestErr)

		return
	}

	ctx.JSON(http.StatusOK, leaveRequest)
}

func (h *LeaveRequest) CreateLeaveRequest(ctx *gin.Context) {
	var createLeaveRequestRequest dto.CreateLeaveRequestRequest
	userIDRaw, _ := ctx.Get("userId")
	userID := userIDRaw.(int)

	if err := util.StrictBindJSON(ctx, &createLeaveRequestRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.New().Struct(createLeaveRequestRequest); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[fe.Field()] = util.MsgForTag(fe)
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if createLeaveRequestRequest.StartDate.After(createLeaveRequestRequest.EndDate) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errors": map[string]string{"startDate": "startDate cannot be after endDate"},
		})
		return
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	start := createLeaveRequestRequest.StartDate

	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)

	if startDate.Before(nowDate) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errors": map[string]string{"startDate": "Leave request cannot be in the past"},
		})
		return
	}

	createLeaveRequestResponse, signupError := h.leaveRequestUsecase.CreateLeaveRequest(&createLeaveRequestRequest, userID)
	if signupError != nil {
		ctx.AbortWithStatusJSON(signupError.Code, signupError)

		return
	}

	ctx.JSON(http.StatusCreated, createLeaveRequestResponse)
}

func (h *LeaveRequest) Approve(ctx *gin.Context) {
	leaveRequestID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Leave Request ID not valid"})

		return
	}

	approveError := h.leaveRequestUsecase.Approve(leaveRequestID)
	if approveError != nil {
		ctx.AbortWithStatusJSON(approveError.Code, approveError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Leave Request Approved"})
}

func (h *LeaveRequest) Reject(ctx *gin.Context) {
	leaveRequestID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Leave Request ID not valid"})

		return
	}

	rejectError := h.leaveRequestUsecase.Reject(leaveRequestID)
	if rejectError != nil {
		ctx.AbortWithStatusJSON(rejectError.Code, rejectError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Leave Request Rejected"})
}
