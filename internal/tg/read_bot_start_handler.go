package tg

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	newsmodels "gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/repositories"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/services"
	"time"
)

func ReadBotStartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	username := update.Message.From.Username
	chatId := update.Message.Chat.ID

	ur := repositories.GetUserRepository()
	user, err := ur.GetUserByTelegramID(userId)
	if err != nil {
		fmt.Printf("start handler failed: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   services.Translate("Error, try again later"),
		})
		if err != nil {
			fmt.Printf("start handler failed: %v", err)
		}
	}

	if user == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		user, err = ur.Create(ctx, userId, username, newsmodels.USER_ROLE_CLIENT)
		if err != nil {
			fmt.Printf("start handler failed: %v", err)
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatId,
				Text:   services.Translate("Error, try again later"),
			})
			if err != nil {
				fmt.Printf("start handler failed to send error-message: %v", err)
			}
		}
	}

	ucr := repositories.GetUserChatRepository()
	uc, err := ucr.GetUserChat(user.ID)
	if err != nil {
		fmt.Printf("start handler failed to retrieve user-chat record from db: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   services.Translate("Error, try again later"),
		})
		if err != nil {
			fmt.Printf("start handler failed to send error-message: %v", err)
		}
	}

	if uc == nil {
		uc, err = ucr.CreateUserChat(ctx, user.ID, chatId)
	}

	// @TODO: Send up to 5 last posts
}
