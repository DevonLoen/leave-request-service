package middleware

import (
	"net/http"
	"strings"

	"github.com/devonLoen/leave-request-service/internal/app/rest_api/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization format must be Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := util.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token: " + err.Error()})
			c.Abort()
			return
		}

		if userIDFloat, ok := claims["id"].(float64); ok {
			c.Set("userId", int(userIDFloat))
		} else {
			c.Set("userId", 0)
		}

		c.Set("authClaims", claims)

		c.Next()
	}
}

func AdminGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsRaw, exists := c.Get("authClaims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
			c.Abort()
			return
		}

		claims, ok := claimsRaw.(map[string]interface{})
		if !ok {
			if mClaims, ok := claimsRaw.(jwt.MapClaims); ok {
				claims = map[string]interface{}(mClaims)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
				c.Abort()
				return
			}
		}

		if role, ok := claims["role"].(string); !ok || (role != "superadmin" && role != "admin") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: admin only"})
			c.Abort()
			return
		}

		c.Next()
	}
}
