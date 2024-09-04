package forum

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/LLIEPJIOK/forum/internal/controller"
	"github.com/LLIEPJIOK/forum/internal/database"
	"github.com/LLIEPJIOK/forum/internal/router"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	forumLogsFile = "logs/forum.txt"
	address       = "localhost:8080"
)

func Start() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("cannot open db connection: %w", err)
	}

	db := database.New(gormDB)
	if err := db.Migrate(); err != nil {
		return fmt.Errorf("cannot up migrations: %w", err)
	}

	file, err := os.OpenFile(forumLogsFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("cannot open file %q: %w", forumLogsFile, err)
	}

	logger := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{}))
	ctrl := controller.New(db, logger)

	rout := router.New(ctrl)
	rout.Run(address)

	return nil
}
