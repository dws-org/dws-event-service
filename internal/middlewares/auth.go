package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/oskargbc/dws-event-service.git/configs"
	"github.com/oskargbc/dws-event-service.git/internal/pkg/supabase"

	"github.com/gin-gonic/gin"
)

// UserContextKey is the key used to store user information in the context
type UserContextKey string

const (
	UserIDKey UserContextKey = "user_id"
)

// AuthMiddleware extracts user ID from Supabase JWT token and stores it in context
// Debugs headers, token, and Supabase response
func AuthMiddleware() gin.HandlerFunc {
	config := configs.GetEnvConfig()
	supabaseClient := supabase.NewSupabaseClient(config)

	return func(c *gin.Context) {
		// Debug: Print all headers

		tokenStr := c.GetHeader("Authorization")

		if tokenStr == "" {
			fmt.Fprintln(os.Stderr, "No Authorization header found")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Authorization header required",
			})
			return
		}

		parts := strings.Split(tokenStr, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Fprintf(os.Stderr, "Invalid Authorization format: %v\n", parts)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid authorization format",
			})
			return
		}

		token := parts[1]

		// Verify the Supabase JWT token
		userID, err := supabaseClient.VerifyJWT(token)
		if err != nil {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid token: " + err.Error(),
			})
			return
		}

		// Store user ID in context
		ctx := context.WithValue(c.Request.Context(), UserIDKey, userID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetUserIDFromContext extracts user ID from the request context
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Request.Context().Value(UserIDKey).(string)
	return userID, exists
}
