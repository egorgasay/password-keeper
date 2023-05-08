package main

import (
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"password-keeper/config"
	"password-keeper/internal/bot"
	"password-keeper/internal/storage"
	"password-keeper/internal/usecase"
	"syscall"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	store, err := storage.New()
	if err != nil {
		log.Fatalf("storage error: %s", err)
	}

	logic := usecase.New(store)

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("zap error: %s", err)
	}

	b, err := bot.New(cfg.Token, cfg.EncryptionKey, logic, logger)
	if err != nil {
		log.Fatalf("bot error: %s", err)
	}

	go func() {
		err := b.Start()
		if err != nil {
			log.Fatalf("bot error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutdown Server ...")

	b.Stop()
}
