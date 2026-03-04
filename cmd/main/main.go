package main

import (
	"log/slog"
	"os"
	"pizza-tracker-go/internal/models"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := loadConfig()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbModel, err := models.InitDB(cfg.DBPath)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	slog.Info("Database initialized successfully")

	sessionStore := setupSessionStore(dbModel.DB, []byte(cfg.SessionSecretKey))

	RegisterCustomValidators()

	h := NewHandler(dbModel)

	router := gin.Default()

	if err := loadTemplates(router); err != nil {
		slog.Error("Failed to load templates", "error", err)
		os.Exit(1)
	}

	setupRoutes(router, h, sessionStore)

	slog.Info("Server running on url", "url", "http://localhost:"+cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		slog.Error("Failed to run server", "error", err)
		os.Exit(1)
	}
}
