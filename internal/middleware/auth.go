package middleware

import (
	"net/http"
	"strings"

	"tienda-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService ports.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		userID, role, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in token"})
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user role"})
			c.Abort()
			return
		}

		isAllowed := false
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
