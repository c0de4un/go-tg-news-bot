package tg

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/repositories"
)

func onInlineKeyboardSelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	cmdKey := string(data)

	userId := mes.Message.From.ID
	chatId := mes.Message.Chat.ID
	ur := repositories.GetUserRepository()
	user, err := ur.GetUserByTelegramID(userId)
	if err != nil {
		fmt.Printf("onInlineKeyboardSelect: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "Error, try again later",
		})
		if err != nil {
			fmt.Printf("onInlineKeyboardSelect: %v", err)
		}
	}

	if cmdKey == "1-1" {
		handlePostCreateStart(user, ctx, b, mes)
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Message.Chat.ID,
		Text:   "You selected: " + cmdKey,
	})
	if err != nil {
		fmt.Printf("onInlineKeyboardSelect: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   "Failed, try again later",
		})
		if err != nil {
			fmt.Printf("onInlineKeyboardSelect: %v", err)
		}
	}
}
