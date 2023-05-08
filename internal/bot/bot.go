package bot

import (
	"fmt"
	"go.uber.org/zap"
	"password-keeper/internal/usecase"
	"time"

	// telegram SDK
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot represents a bot.
type Bot struct {
	token  string
	logic  *usecase.UseCase
	logger *zap.Logger

	*tgapi.BotAPI

	stopHiding   func()
	toHide       chan messageInfo
	hideInterval int64

	langStorage map[string]messages
}

type messages struct {
	Russian string
	English string
}

type keyboards struct {
	Russian tgapi.InlineKeyboardMarkup
	English tgapi.InlineKeyboardMarkup
}

// New creates a new bot.
func New(token string, deletionInterval time.Duration, logic *usecase.UseCase, logger *zap.Logger) (*Bot, error) {
	bot, err := tgapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("error creating bot: %w", err)
	}

	return &Bot{
		token:        token,
		logic:        logic,
		BotAPI:       bot,
		logger:       logger,
		hideInterval: int64(deletionInterval.Seconds()),
	}, nil
}

// Start starts the bot.
func (b *Bot) Start() {
	u := tgapi.NewUpdate(0)
	u.Timeout = 60
	b.toHide, b.stopHiding = b.Watch()

	//startMessage = fmt.Sprintf(startMessage, b.hideInterval)

	updates := b.GetUpdatesChan(u)
	for update := range updates {
		if update.CallbackQuery != nil {
			b.handleCallbackQuery(update.CallbackQuery)
			continue
		}

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		}

		b.handleMessage()
	}

}

// Stop stops the bot.
func (b *Bot) Stop() {
	b.StopReceivingUpdates()
	b.stopHiding()
}
