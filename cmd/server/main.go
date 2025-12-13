package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/deeploy-sh/deeploy/internal/server/app"
	"github.com/deeploy-sh/deeploy/internal/server/config"
	"github.com/deeploy-sh/deeploy/internal/server/logger"
	"github.com/deeploy-sh/deeploy/internal/server/routes"
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
