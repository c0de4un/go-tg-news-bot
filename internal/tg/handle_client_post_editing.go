package tg

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	newsmodels "gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/repositories"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/services"
)

func handlePostCreateStart(user *newsmodels.UserModel, uc *newsmodels.ChatModel, c *newsmodels.ClientModel, ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	ucr := repositories.GetUserChatRepository()
	err := ucr.SetState(ctx, uc.ID, newsmodels.CHAT_STATE_FORWARDED_POST_INPUT)
	if err != nil {
		fmt.Printf("handlePostCreateStart: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   services.Translate("Failed, try again later"),
		})
		if err != nil {
			fmt.Printf("handlePostCreateStart: %v", err)
		}

		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Message.Chat.ID,
		Text:   services.Translate("Enter post content. Images and other media are supported"),
	})
	if err != nil {
		fmt.Printf("handlePostCreateStart: %v", err)
	}
}
