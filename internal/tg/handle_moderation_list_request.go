package tg

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	newsmodels "gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/repositories"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/services"
)

func handleModerationListRequest(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	fr := repositories.GetForwardPostRepository()

	post, err := fr.GetByStatus(newsmodels.POST_STATUS_MODERATION)
	if err != nil {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   services.Translate("Error, try again later"),
		})
		if err != nil {
			fmt.Printf("\nhandleModerationListRequest: failed to retrieve menu for admin: %v\n", err)
		}

		return
	}
	if post == nil {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   services.Translate("No posts to moderate"),
		})

		return
	}

	ts := services.GetTelegramService()
	err = ts.SendForwardPost(post, mes.Message.Chat.ID, false)
	if err != nil {
		fmt.Printf("handleModerationListRequest: failed to render post for admin: %v", err)
	}

	kb := inline.New(b).
		Row().
		Button(services.Translate("Publish"), []byte(fmt.Sprintf("2-3;%d", post.ID)), onInlineKeyboardSelect)
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Message.Chat.ID,
		Text:        services.Translate("Menu"),
		ReplyMarkup: kb,
	})
	if err != nil {
		fmt.Printf("handleModerationListRequest: failed to render moderation menu for admin: %v", err)
	}
}
