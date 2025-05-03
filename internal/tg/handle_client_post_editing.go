package tg

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	newsmodels "gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/repositories"
)

func handlePostCreateStart(user *newsmodels.UserModel, uc *newsmodels.ChatModel, c *newsmodels.ClientModel, ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	// @TODO: handlePostCreateStart()

	ucr := repositories.GetUserChatRepository()
	err := ucr.SetState(ctx, uc.ID, newsmodels.CHAT_STATE_POST_TITLE)
	if err != nil {
		fmt.Printf("handlePostCreateStart: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Message.Chat.ID,
			Text:   "Failed, try again later",
		})
		if err != nil {
			fmt.Printf("handlePostCreateStart: %v", err)
		}

		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Message.Chat.ID,
		Text:   "Please, enter post title",
	})
	if err != nil {
		fmt.Printf("handlePostCreateStart: %v", err)
	}
}
