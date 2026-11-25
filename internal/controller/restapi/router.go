package restapi

import (
	"net/http"

	"bukus/config"
	"bukus/internal/controller/restapi/middleware"
	v1 "bukus/internal/controller/restapi/v1"
	"bukus/internal/usecase"
	"bukus/pkg/jwt"
	"bukus/pkg/logger"

	"github.com/gofiber/fiber/v3"
)

// NewRouter -.
func NewRouter(app *fiber.App, cfg *config.Config, b usecase.Book, a usecase.Auth, jwtAccess *jwt.Manager, l logger.Interface) {
	// Middleware
	app.Use(middleware.CORS())
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	// K8s probe
	app.Get("/healthz", func(ctx fiber.Ctx) error {
		return ctx.SendStatus(http.StatusOK)
	})

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewRoutes(apiV1Group, b, a, jwtAccess, l)
	}
}
