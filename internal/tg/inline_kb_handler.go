package tg

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	newsmodels "gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/repositories"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/services"
	"strconv"
	"strings"
)

func onInlineKeyboardSelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	dataParts := strings.Split(string(data), ";")
	cmdKey := dataParts[0]

	chatId := mes.Message.Chat.ID
	userId, err := strconv.ParseInt(dataParts[1], 10, 64)
	if err != nil {
		fmt.Printf("onInlineKeyboardSelect: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   services.Translate("Failed, try again later"),
		})
		if err != nil {
			fmt.Printf("onInlineKeyboardSelect: %v", err)
		}
	}

	ur := repositories.GetUserRepository()
	user, err := ur.GetUserByTelegramID(userId)
	if err != nil {
		fmt.Printf("onInlineKeyboardSelect: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   services.Translate("Failed, try again later"),
		})
		if err != nil {
			fmt.Printf("onInlineKeyboardSelect: %v", err)
		}
	}

	if user == nil {
		fmt.Printf("onInlineKeyboardSelect: user not found for tg-id %v", userId)
		return
	}

	cr := repositories.GetClientRepository()
	c, err := cr.GetClientByUserID(user.ID)
	if err != nil {
		fmt.Printf("start handler failed to retrieve client from db: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   services.Translate("Failed, try again later"),
		})
		if err != nil {
			fmt.Printf("start handler failed to send error-message: %v", err)
		}
	}

	ucr := repositories.GetUserChatRepository()
	uc, err := ucr.GetUserChat(user.ID, newsmodels.CHAT_TYPE_EDITOR, services.GetEditBotID())
	if err != nil {
		fmt.Printf("start handler failed to retrieve user-chat record from db: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   services.Translate("Failed, try again later"),
		})
		if err != nil {
			fmt.Printf("start handler failed to send error-message: %v", err)
		}
	}

	if cmdKey == "1-1" {
		handlePostCreateStart(user, uc, c, ctx, b, mes)
		return
	}

	if err != nil {
		fmt.Printf("onInlineKeyboardSelect: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   services.Translate("Failed, try again later"),
		})
		if err != nil {
			fmt.Printf("onInlineKeyboardSelect: %v", err)
		}
	}
}
