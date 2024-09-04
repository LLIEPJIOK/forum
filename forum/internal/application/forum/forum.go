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

func Start() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("cannot open db connection: %w", err)
	}

	db := database.New(gormDB)
	if err := db.Migrate(); err != nil {
		return fmt.Errorf("cannot up migrations: %w", err)
	}

	file, err := os.OpenFile(os.Getenv("LOGS_FILE"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("cannot open file %q: %w", os.Getenv("LOGS_FILE"), err)
	}

	logger := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{}))
	ctrl := controller.New(db, logger)

	rout := router.New(ctrl)
	rout.Run(os.Getenv("API_ADDRESS"))

	return nil
}
