package bot

import (
	"errors"
	"fmt"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"password-keeper/internal/storage"
	"strings"
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

// handleStart handles start command.
func (b *Bot) handleStart(msg *tgapi.Message) {
	msgConfig := tgapi.NewMessage(msg.Chat.ID, startMessage)

	_, err := b.Send(msgConfig)
	if err != nil {
		log.Println("send error: ", err)
	}
}

func (b *Bot) handleSet(msg *tgapi.Message) {
	split := strings.Split(msg.Text, " ")

	msgConfig := tgapi.NewMessage(msg.Chat.ID, setMessage)
	if len(split) != 4 {
		msgConfig = tgapi.NewMessage(msg.Chat.ID, "Invalid command")
	}

	err := b.logic.Save(msg.Chat.ID, split[1], split[2], split[3])
	if err != nil {
		msgConfig.Text = "Сохранение не удалась"
		log.Printf("save error: %v\n", err)
	}

	_, err = b.Send(msgConfig)
	if err != nil {
		log.Println("send error: ", err)
	}
}

func (b *Bot) handleGet(msg *tgapi.Message) {
	split := strings.Split(msg.Text, " ")

	msgConfig := tgapi.NewMessage(msg.Chat.ID, setMessage)
	if len(split) != 2 {
		msgConfig = tgapi.NewMessage(msg.Chat.ID, "Invalid command")
	}
	service := split[1]

	pair, err := b.logic.Get(msg.Chat.ID, service)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			msgConfig.Text = "Сервис не найден"
		} else {
			msgConfig.Text = "Сохранение не удалась"
		}
		log.Printf("save error: %v\n", err)
	} else {
		msgConfig.Text = fmt.Sprintf("%s \nЛогин:  %s\nПароль: %s", service, pair.Login, pair.Password)
	}

	_, err = b.Send(msgConfig)
	if err != nil {
		log.Println("send error: ", err)
	}
}

func (b *Bot) handleDel(msg *tgapi.Message) {
	split := strings.Split(msg.Text, " ")

	msgConfig := tgapi.NewMessage(msg.Chat.ID, setMessage)
	if len(split) != 2 {
		msgConfig = tgapi.NewMessage(msg.Chat.ID, "Invalid command")
	}
	service := split[1]

	err := b.logic.Delete(msg.Chat.ID, service)
	if err != nil {
		msgConfig.Text = "Сохранение не удалась"
		log.Printf("save error: %v\n", err)
	} else {
		msgConfig.Text = fmt.Sprintf("Удалено! :)")
	}

	_, err = b.Send(msgConfig)
	if err != nil {
		log.Println("send error: ", err)
	}
}
