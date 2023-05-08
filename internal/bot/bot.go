package bot

import (
	"fmt"
	"go.uber.org/zap"
	"password-keeper/internal/usecase"

	// telegram SDK
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot represents a bot.
type Bot struct {
	token  string
	logic  *usecase.UseCase
	logger *zap.Logger

	*tgapi.BotAPI
}

// New creates a new bot.
func New(token string, key string, logic *usecase.UseCase, logger *zap.Logger) (*Bot, error) {
	bot, err := tgapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("error creating bot: %w", err)
	}

	return &Bot{
		token:  token,
		logic:  logic,
		BotAPI: bot,
		logger: logger,
	}, nil
}

// Start starts the bot.
func (b *Bot) Start() error {
	u := tgapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		}

		b.handleMessage()
	}

	return nil
}

// Stop stops the bot.
func (b *Bot) Stop() {
	b.StopReceivingUpdates()
}
