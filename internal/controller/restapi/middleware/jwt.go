package middleware

import (
	"bukus/pkg/jwt"
	"bukus/pkg/logger"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// JWTMiddleware validates the access token from the Authorization header.
func JWTMiddleware(jwtManager *jwt.Manager, l logger.Interface) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return fmt.Errorf("missing authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fmt.Errorf("missing authorization header")
		}

		userID, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			l.Error(err, "restapi - v1 - JWTMiddleware - ParseToken")
			return fmt.Errorf("missing authorization header")
		}

		ctx.Locals("userID", userID)

		return ctx.Next()
	}
}
