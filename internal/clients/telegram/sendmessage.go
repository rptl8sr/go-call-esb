package telegram

import (
	"context"
	"fmt"
)

var (
	sendMessagePath = "/sendMessage"
)

type SendMessageRequest struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

func (c *Client) SendMessage(ctx context.Context, text string) error {
	body := &SendMessageRequest{
		ChatID: c.chatID,
		Text:   text,
	}

	_, err := c.doPost(ctx, sendMessagePath, body)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
