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
	pr := repositories.GetPostRepository()
	p, err := pr.GetPostByStatus(newsmodels.POST_STATUS_MODERATION)
	if err != nil {
		fmt.Printf("handleModerationListRequest: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   services.Translate("Failed, try again later"),
		})
		if err != nil {
			fmt.Printf("handleModerationListRequest: %v", err)
		}
		return
	}

	ts := services.GetTelegramService()
	err = ts.SendPost(ctx, p, mes.Message.Chat.ID)
	if err != nil {
		fmt.Printf("handleModerationListRequest: failed to render post for admin: %v", err)
	}

	kb := inline.New(b).
		Row().
		Button(services.Translate("Publish"), []byte(fmt.Sprintf("2-3;%d", p.ID)), onInlineKeyboardSelect).
		Button(services.Translate("Delete"), []byte(fmt.Sprintf("2-4;%d", p.ID)), onInlineKeyboardSelect).
		Button(services.Translate("Next"), []byte("2-5"), onInlineKeyboardSelect)
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Message.Chat.ID,
		Text:        services.Translate("Menu"),
		ReplyMarkup: kb,
	})
	if err != nil {
		fmt.Printf("handleModerationListRequest: failed to render moderation menu for admin: %v", err)
	}
}
