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
	"strconv"
)

func handleModerationPublishRequest(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, dataParts []string) {
	postId, err := strconv.Atoi(dataParts[1])
	if err != nil {
		fmt.Printf("error at handleModerationPublishRequest: %v", err)
		return
	}

	fr := repositories.GetForwardPostRepository()
	post, err := fr.GetById(int64(postId))
	if err != nil {
		fmt.Printf("handleModerationPublishRequest: failed to retrieve post by id: %v", err)
		return
	}

	post.Status = newsmodels.POST_STATUS_PUBLISHED
	err = fr.SetState(post.ID, newsmodels.POST_STATUS_PUBLISHED)
	if err != nil {
		fmt.Printf("handleModerationPublishRequest: failed to set post status: %v", err)
		return
	}

	ts := services.GetTelegramService()

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Message.Chat.ID,
		Text:   services.Translate("Processing . . ."),
	})
	if err != nil {
		fmt.Printf("handleModerationPublishRequest: %v", err)
	}

	ts.PublishForwardPost(post)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Message.Chat.ID,
		Text:   services.Translate("Completed"),
	})
	if err != nil {
		fmt.Printf("handleModerationPublishRequest: %v", err)
	}

	ur := repositories.GetUserRepository()
	user, err := ur.GetUserByTelegramID(mes.Message.From.ID)
	if err != nil {
		fmt.Printf("handleModerationPublishRequest: %v", err)
		return
	}

	if user != nil {
		kb := inline.New(b).
			Row().
			Button(services.Translate("Create post"), []byte(fmt.Sprintf("1-1;%d", user.ID)), onInlineKeyboardSelect).
			Button(services.Translate("Moderate"), []byte(fmt.Sprintf("2-1;%d", user.ID)), onInlineKeyboardSelect)

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      mes.Message.Chat.ID,
			Text:        services.Translate("Menu"),
			ReplyMarkup: kb,
		})
	}
}
