package v1

import (
	"bukus/internal/controller/restapi/middleware"
	"bukus/internal/usecase"
	"bukus/pkg/jwt"
	"bukus/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

// NewRoutes -.
func NewRoutes(apiV1Group fiber.Router, b usecase.Book, a usecase.Auth, jwtManager *jwt.Manager, l logger.Interface) {
	r := &V1{
		bookUC: b,
		authUC: a,
		l:      l,
		v:      validator.New(validator.WithRequiredStructEnabled()),
	}

	// Auth routes (public)
	authGroup := apiV1Group.Group("/auth")
	{
		authGroup.Post("/register", r.register)
		authGroup.Post("/login", r.login)
		authGroup.Post("/refresh", r.refresh)
	}

	// Book routes (protected)
	bookGroup := apiV1Group.Group("/books", middleware.JWTMiddleware(jwtManager, l))
	{
		bookGroup.Get("/", r.list)
		bookGroup.Post("/", r.store)
		bookGroup.Get("/:id", r.getByID)
		bookGroup.Put("/:id", r.update)
		bookGroup.Delete("/:id", r.delete)
	}
}
