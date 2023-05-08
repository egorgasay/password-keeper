package bot

import (
	"errors"
	"fmt"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"password-keeper/internal/storage"
	"strings"
	"time"
)

// handleMessage handles commands.
func (b *Bot) handleCommand(msg *tgapi.Message) {
	switch msg.Command() {
	case start:
		b.handleStart(msg)
	case set:
		b.handleSet(msg)
	case get:
		b.handleGet(msg)
	case del:
		b.handleDel(msg)
	}
}

// handleMessage handles messages.
func (b *Bot) handleMessage() {

}

func (b *Bot) handleMessageLang(msg string, chatID int64) string {
	lang := b.logic.GetLang(chatID)
	switch lang {
	case "ru":
		return allMessages[msg].Russian
	default:
		return allMessages[msg].English
	}
}

func (b *Bot) handleKeyboardLang(keyboard string, chatID int64) tgapi.InlineKeyboardMarkup {
	lang := b.logic.GetLang(chatID)
	switch lang {
	case "ru":
		return allKeyboards[keyboard].Russian
	default:
		return allKeyboards[keyboard].English
	}
}

// handleStart handles start command.
func (b *Bot) handleStart(msg *tgapi.Message) {
	msgConfig := tgapi.NewMessage(msg.Chat.ID,
		fmt.Sprintf(b.handleMessageLang(start, msg.Chat.ID),
			b.hideInterval))
	msgConfig.ReplyMarkup = b.handleKeyboardLang(startKeyboard, msg.Chat.ID)

	_, err := b.Send(msgConfig)
	if err != nil {
		log.Println("send error: ", err)
	}
}

func (b *Bot) handleSet(msg *tgapi.Message) {
	split := strings.Split(msg.Text, " ")

	msgConfig := tgapi.NewMessage(msg.Chat.ID, b.handleMessageLang(set, msg.Chat.ID))
	if len(split) != 4 {
		msgConfig = tgapi.NewMessage(msg.Chat.ID, b.handleMessageLang(wrongInputErr, msg.Chat.ID))
		m, err := b.Send(msgConfig)
		if err != nil {
			log.Println("send error: ", err)
		} else {
			b.toHide <- messageInfo{
				chatID:    msg.Chat.ID,
				id:        msg.MessageID,
				createdAt: time.Now(),
			}

			b.toHide <- messageInfo{
				chatID:    m.Chat.ID,
				id:        m.MessageID,
				createdAt: time.Now(),
			}
		}
		return
	}

	err := b.logic.Save(msg.Chat.ID, split[1], split[2], split[3])
	if err != nil {
		msgConfig.Text = b.handleMessageLang(setErr, msg.Chat.ID)
		log.Printf("save error: %v\n", err)
	}

	m, err := b.Send(msgConfig)
	if err != nil {
		log.Println("send error: ", err)
	} else {
		b.toHide <- messageInfo{
			chatID:    msg.Chat.ID,
			id:        msg.MessageID,
			createdAt: time.Now(),
		}

		b.toHide <- messageInfo{
			chatID:    m.Chat.ID,
			id:        m.MessageID,
			createdAt: time.Now(),
		}
	}
}

func (b *Bot) handleGet(msg *tgapi.Message) {
	split := strings.Split(msg.Text, " ")

	msgConfig := tgapi.NewMessage(msg.Chat.ID, b.handleMessageLang(get, msg.Chat.ID))
	if len(split) != 2 {
		msgConfig = tgapi.NewMessage(msg.Chat.ID, b.handleMessageLang(wrongInputErr, msg.Chat.ID))
		m, err := b.Send(msgConfig)
		if err != nil {
			log.Println("send error: ", err)
		} else {
			b.toHide <- messageInfo{
				chatID:    msg.Chat.ID,
				id:        msg.MessageID,
				createdAt: time.Now(),
			}

			b.toHide <- messageInfo{
				chatID:    m.Chat.ID,
				id:        m.MessageID,
				createdAt: time.Now(),
			}
		}
		return
	}
	service := split[1]

	pair, err := b.logic.Get(msg.Chat.ID, service)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			msgConfig.Text = b.handleMessageLang(serviceNotFoundErr, msg.Chat.ID)
		} else {
			msgConfig.Text = b.handleMessageLang(getErr, msg.Chat.ID)
		}
		log.Printf("get error: %v\n", err)
	} else {
		msgConfig.ReplyMarkup = b.handleKeyboardLang(hideKeyboard, msg.Chat.ID)
		msgConfig.Text = fmt.Sprintf(b.handleMessageLang(get, msg.Chat.ID), service, pair.Login, pair.Password)
	}

	m, err := b.Send(msgConfig)
	if err != nil {
		log.Println("send error: ", err)
	} else {
		b.toHide <- messageInfo{
			chatID:    msg.Chat.ID,
			id:        msg.MessageID,
			createdAt: time.Now(),
		}

		b.toHide <- messageInfo{
			chatID:    m.Chat.ID,
			id:        m.MessageID,
			createdAt: time.Now(),
		}
	}
}

func (b *Bot) handleDel(msg *tgapi.Message) {
	split := strings.Split(msg.Text, " ")

	msgConfig := tgapi.NewMessage(msg.Chat.ID, b.handleMessageLang(del, msg.Chat.ID))
	if len(split) != 2 {
		msgConfig = tgapi.NewMessage(msg.Chat.ID, b.handleMessageLang(wrongInputErr, msg.Chat.ID))
	}
	service := split[1]

	err := b.logic.Delete(msg.Chat.ID, service)
	if err != nil {
		msgConfig.Text = b.handleMessageLang(delErr, msg.Chat.ID)
		log.Printf("del error: %v\n", err)
	}

	m, err := b.Send(msgConfig)
	if err != nil {
		log.Println("send error: ", err)
	} else {
		b.toHide <- messageInfo{
			chatID:    msg.Chat.ID,
			id:        msg.MessageID,
			createdAt: time.Now(),
		}

		b.toHide <- messageInfo{
			chatID:    m.Chat.ID,
			id:        m.MessageID,
			createdAt: time.Now(),
		}
	}
}

// handleMessage handle callbacks from user.
func (b *Bot) handleCallbackQuery(query *tgapi.CallbackQuery) {
	split := strings.Split(query.Data, "::")
	if len(split) == 0 {
		return
	}

	defer b.logger.Sync()

	text := split[0]

	switch text {
	case hide:
		msg := tgapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
		if _, err := b.Send(msg); err != nil {
			b.logger.Warn(fmt.Sprintf("del error: %v", err.Error()))
		}

		msg = tgapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID-1)
		if _, err := b.Send(msg); err != nil {
			b.logger.Warn(fmt.Sprintf("del error: %v", err.Error()))
		}
	case changeLang:
		msg := tgapi.NewEditMessageTextAndMarkup(
			query.Message.Chat.ID,
			query.Message.MessageID,
			"Choose a new language 🌎",
			b.handleKeyboardLang(setLangKeyboard, query.Message.Chat.ID),
		)

		if _, err := b.Send(msg); err != nil {
			b.logger.Warn(fmt.Sprintf("send error: %v", err.Error()))
		}
	case change:
		if len(split) == 1 {
			return
		}

		b.logic.SetLang(query.Message.Chat.ID, split[1])

		msg := tgapi.NewEditMessageTextAndMarkup(
			query.Message.Chat.ID, query.Message.MessageID,
			fmt.Sprintf(b.handleMessageLang(start, query.Message.Chat.ID), b.hideInterval),
			b.handleKeyboardLang(startKeyboard, query.Message.Chat.ID),
		)

		if _, err := b.Send(msg); err != nil {
			b.logger.Warn(fmt.Sprintf("send error: %v", err.Error()))
		}

	}
}
