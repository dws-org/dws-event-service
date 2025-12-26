package router

import (
	"github.com/dws-org/dws-event-service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware middleware for protected routes requiring user authentication
func AuthMiddleware() gin.HandlerFunc {
	return middlewares.AuthMiddleware()
}

// AdminAuthMiddleware middleware for admin-only routes
// Since role-based access is not needed, this just checks authentication
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Just check if user is authenticated
		AuthMiddleware()(c)

		// If authentication passed, continue
		if c.IsAborted() {
			return
		}

		c.Next()
	}
}
