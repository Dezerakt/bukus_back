// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"bukus/config"
	"bukus/internal/controller/restapi"
	persistBookRepo "bukus/internal/repo/persistent/book"
	persistUserRepo "bukus/internal/repo/persistent/user"
	"bukus/internal/usecase"
	authUC "bukus/internal/usecase/auth"
	bookUC "bukus/internal/usecase/book"
	"bukus/pkg/httpserver"
	"bukus/pkg/jwt"
	"bukus/pkg/logger"
	"bukus/pkg/postgres"
)

type useCases struct {
	book usecase.Book
	auth usecase.Auth
}

type servers struct {
	http *httpserver.Server
}

func initUseCases(cfg *config.Config, pg *postgres.Postgres) useCases {
	bookRepo := persistBookRepo.New(pg)
	userRepo := persistUserRepo.New(pg)

	jwtAccess := jwt.New(cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry)
	jwtRefresh := jwt.New(cfg.JWT.RefreshSecret, cfg.JWT.RefreshTokenExpiry)

	return useCases{
		book: bookUC.New(bookRepo),
		auth: authUC.New(userRepo, jwtAccess, jwtRefresh),
	}
}

func initServers(cfg *config.Config, uc useCases, l logger.Interface) servers {
	jwtAccess := jwt.New(cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry)

	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
	restapi.NewRouter(httpServer.App, cfg, uc.book, uc.auth, jwtAccess, l)

	return servers{
		http: httpServer,
	}
}

func (obj *servers) startServers() {
	obj.http.Start()
}

func (obj *servers) waitForShutdown(l logger.Interface) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Info("app - Run - signal: %s", sig.String())
	case err = <-obj.http.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	obj.shutdownServers(l)
}

func (obj *servers) shutdownServers(l logger.Interface) {
	if err := obj.http.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	uc := initUseCases(cfg, pg)
	s := initServers(cfg, uc, l)
	s.startServers()
	s.waitForShutdown(l)
}
