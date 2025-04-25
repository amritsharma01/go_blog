package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"crud_api/models"
	"crud_api/utils"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// JWTMiddlewareConfig holds the configuration for the JWT middleware
type JWTMiddlewareConfig struct {
	DB *gorm.DB
}

// NewJWTMiddleware creates a new instance of the JWT middleware with dependencies
func NewJWTMiddleware(db *gorm.DB) *JWTMiddlewareConfig {
	return &JWTMiddlewareConfig{DB: db}
}

// Middleware is a middleware to protect routes by verifying JWT tokens.
func (config *JWTMiddlewareConfig) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponse(c, http.StatusUnauthorized, "Missing authorization header")
		}

		// Split the header into "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
		}
		tokenString := parts[1]

		// Parse and verify the token
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token's signing method is valid (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil {
			return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
		}

		// Verify the token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims")
		}

		// Check token expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Token has expired")
			}
		}

		// Retrieve the user ID from the claims
		var userID uint
		switch v := claims["user_id"].(type) {
		case float64:
			userID = uint(v) // Handle JSON number case
		case string:
			id, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID format")
			}
			userID = uint(id)
		default:
			return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID type in token")
		}

		// Find the user in database
		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			return utils.ErrorResponse(c, http.StatusUnauthorized, "User not found")
		}

		// Set the user object into the context
		c.Set("user", user)

		// Continue with the next handler
		return next(c)
	}
}
