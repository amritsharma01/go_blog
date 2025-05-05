package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	responsemodels "crud_api/response_models"
	"crud_api/services"

	"os"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JWTMiddlewareConfig struct {
	UserService services.UserService
}

func NewJWTMiddleware(userService services.UserService) *JWTMiddlewareConfig {
	return &JWTMiddlewareConfig{UserService: userService}
}

func (config *JWTMiddlewareConfig) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return responsemodels.ErrorResponse(c, http.StatusUnauthorized, "Missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return responsemodels.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
		}
		tokenString := parts[1]

		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil {
			return responsemodels.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return responsemodels.ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims")
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return responsemodels.ErrorResponse(c, http.StatusUnauthorized, "Token has expired")
			}
		}

		var userID uint
		switch v := claims["user_id"].(type) {
		case float64:
			userID = uint(v)
		case string:
			id, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return responsemodels.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID format")
			}
			userID = uint(id)
		default:
			return responsemodels.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID type")
		}

		user, err := config.UserService.GetByID(userID)
		if err != nil {
			return responsemodels.ErrorResponse(c, http.StatusUnauthorized, "User not found")
		}

		c.Set("user", *user)
		return next(c)
	}
}
