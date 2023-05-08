package config

import (
	"errors"
	"flag"
	"os"
	"time"
)

// Flag struct for parsing from env and cmd args.
type Flag struct {
	EncryptionKey    *string
	Token            *string
	DeletionInterval *time.Duration
	Storage          *string
	DSN              *string
}

var (
	f Flag

	// ErrTokenNotSet error when the key is not set.
	ErrTokenNotSet = errors.New("telegram-bot token is not set")

	// ErrEncryptionKeyNotSet error when the key is not set.
	ErrEncryptionKeyNotSet = errors.New("encryption-key is not set")
)

func init() {
	f.EncryptionKey = flag.String("key", "", "-key=KEY")
	f.Token = flag.String("token", "", "-token=TOKEN")
	f.DeletionInterval = flag.Duration("interval", 7*time.Second, "-interval=1s")
	f.Storage = flag.String("storage", "sqlite", "-storage=sqlite|postgres")
	f.DSN = flag.String("dsn", "keeper.db", "-dsn=CONNECTION_STRING")
}

// Config contains all the settings for configuring the application.
type Config struct {
	EncryptionKey    string
	Token            string
	DeletionInterval time.Duration
	Storage          string
	DSN              string
}

// New initializing the config for the application.
func New() (*Config, error) {
	flag.Parse()

	if key, ok := os.LookupEnv("ENCRYPTION_KEY"); ok {
		*f.EncryptionKey = key
	}

	if key, ok := os.LookupEnv("TELEGRAM_API_KEY"); ok {
		*f.Token = key
	}

	if *f.EncryptionKey == "" {
		return nil, ErrEncryptionKeyNotSet
	}

	if *f.Token == "" {
		return nil, ErrTokenNotSet
	}

	return &Config{
		EncryptionKey:    *f.EncryptionKey,
		Token:            *f.Token,
		DeletionInterval: *f.DeletionInterval,
		Storage:          *f.Storage,
		DSN:              *f.DSN,
	}, nil
}
