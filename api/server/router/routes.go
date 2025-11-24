package route

import (
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/handler"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterPublicEndpoints(router *gin.Engine, userHandlers *handler.User, authHandlers *handler.Auth, leaveRequestHandlers *handler.LeaveRequest) {
	public := router.Group("/api/v1")
	{
		public.POST("/auth/login", authHandlers.Login)
	}

	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthGuard())
	{
		protected.POST("/leave-requests", leaveRequestHandlers.CreateLeaveRequest)
		protected.GET("/my-leave-requests", leaveRequestHandlers.GetMyLeaveRequests)
		protected.PATCH("/leave-requests/:id/submit", leaveRequestHandlers.Submit)

	}

	protectedAdmin := router.Group("/api/v1")
	protectedAdmin.Use(middleware.AuthGuard(), middleware.AdminGuard())
	{
		protectedAdmin.GET("/users", userHandlers.GetAllUsers)
		protectedAdmin.GET("/users/:id", userHandlers.GetUser)
		protectedAdmin.POST("/users", userHandlers.CreateUser)

		protectedAdmin.GET("/leave-requests", leaveRequestHandlers.GetAllLeaveRequests)
		protectedAdmin.GET("/leave-requests/:id", leaveRequestHandlers.GetLeaveRequest)
		protectedAdmin.PATCH("/leave-requests/:id/approve", leaveRequestHandlers.Approve)
		protectedAdmin.PATCH("/leave-requests/:id/reject", leaveRequestHandlers.Reject)
	}
}
