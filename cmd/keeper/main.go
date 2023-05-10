package main

import (
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"password-keeper/config"
	"password-keeper/internal/bot"
	"password-keeper/internal/storage"
	"password-keeper/internal/storage/queries"
	"password-keeper/internal/usecase"
	"syscall"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	store, err := storage.New(cfg.Storage, cfg.DSN)
	if err != nil {
		log.Fatalf("storage error: %s", err)
	}

	logic, err := usecase.New(store, cfg.EncryptionKey)
	if err != nil {
		log.Fatalf("logic error: %s", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("zap error: %s", err)
	}

	b, err := bot.New(cfg.Token, cfg.DeletionInterval, logic, logger)
	if err != nil {
		log.Fatalf("bot error: %s", err)
	}

	log.Println("Starting bot...")
	go b.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutdown Server ...")

	b.Stop()
	if err := queries.Close(); err != nil {
		log.Fatalf("queries close error: %s", err)
	}
}
