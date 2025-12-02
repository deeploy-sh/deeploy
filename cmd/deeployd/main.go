package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/deeploy-sh/deeploy/internal/deeployd/app"
	"github.com/deeploy-sh/deeploy/internal/deeployd/config"
	"github.com/deeploy-sh/deeploy/internal/deeployd/logger"
	"github.com/deeploy-sh/deeploy/internal/deeployd/routes"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.IsDevelopment())

	application, err := app.New(cfg)
	if err != nil {
		slog.Error("failed to initialize app", "error", err)
		os.Exit(1)
	}
	defer application.Close()

	handler := routes.Setup(application)
	slog.Info("server starting", "port", cfg.Port)

	err = http.ListenAndServe(":"+cfg.Port, handler)
	if err != nil {
		slog.Error("server failed", "error", err)
	}
}
