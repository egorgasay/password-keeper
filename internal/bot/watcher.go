package bot

import (
	"fmt"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type messageInfo struct {
	chatID    int64
	id        int
	createdAt time.Time
}

func (b *Bot) Watch() (chan messageInfo, func()) {
	messagesCh := make(chan messageInfo, 10000)
	cancelCh := make(chan bool)

	go func() {
		for {
			select {
			case <-cancelCh:
				close(messagesCh)
				return
			case msg := <-messagesCh:
				for time.Now().Unix()-b.hideInterval <= msg.createdAt.Unix() {
				}

				msgDelConfig := tgapi.NewDeleteMessage(msg.chatID, msg.id)
				if _, err := b.Send(msgDelConfig); err != nil {
					b.logger.Warn(fmt.Sprintf("del error: %v", err.Error()))
				}
			}
		}
	}()

	return messagesCh, func() {
		close(cancelCh)
	}
}
