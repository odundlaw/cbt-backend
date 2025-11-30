package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/odundlaw/cbt-backend/internal/config"
)

func main() {
	ctx := context.Background()

	cfg := Config{
		add: ":8080",
		db: DBConfig{
			dsn: config.DatabaseURL,
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.Default()

	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("connected to database", "dsn", cfg.db.dsn)

	api := Application{
		config: cfg,
		conn:   conn,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("server has failed to start", "error", err)
		os.Exit(1)
	}
}
