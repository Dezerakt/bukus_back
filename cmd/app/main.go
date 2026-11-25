package main

import (
	"context"
	"log"
	"os"

	"bukus/config"
	"bukus/internal/app"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	if len(os.Args) > 1 {
		pool, err := pgxpool.New(context.Background(), cfg.PG.URL)
		if err != nil {
			log.Fatalf("main - %s - pgxpool.New: %v", os.Args[1], err)
		}

		switch os.Args[1] {
		case "migrate":
			app.Migrate(pool)
		default:
			log.Fatalf("unknown command: %s", os.Args[1])
		}

		pool.Close()

		return
	}

	app.Run(cfg)
}
