package main

import (
	_log "log"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/gin-router"
	"url-shortener/internal/lib/logger"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		_log.Fatalf("Error on config init: %s", err)
	}

	log := logger.SetupLogger(cfg.Env)

	log.Info("Start url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := gin_router.InitEngine(log, storage)

	router.Run(cfg.HTTPServer.Address)
}
