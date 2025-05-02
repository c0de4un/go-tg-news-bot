package tg

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/repositories"
	"time"
)

func ClientStartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	username := update.Message.From.Username
	chatId := update.Message.Chat.ID

	ur := repositories.GetUserRepository()
	user, err := ur.GetUserByTelegramID(userId)
	if err != nil {
		fmt.Printf("start handler failed: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "Error, try again later",
		})
		if err != nil {
			fmt.Printf("start handler failed: %v", err)
		}
	}
	if user != nil {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   fmt.Sprintf("Welcome back, %s", user.TelegramUsername),
		})
		if err != nil {
			fmt.Printf("start handler failed: %v", err)
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		user, err = ur.Create(ctx, userId, username)
		if err != nil {
			fmt.Printf("start handler failed: %v", err)
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatId,
				Text:   "Failed to register, try again later",
			})
			if err != nil {
				fmt.Printf("start handler failed to send error-message: %v", err)
			}
		}
	}

	ucr := repositories.GetUserChatRepository()
	uc, err := ucr.GetUserChat(user.ID, chatId)
	if err != nil {
		fmt.Printf("start handler failed to retrieve user-chat record from db: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "Failed to register, try again later",
		})
		if err != nil {
			fmt.Printf("start handler failed to send error-message: %v", err)
		}
	}

	if uc == nil {
		uc, err = ucr.CreateUserChat(ctx, user.ID, chatId)
	} else {
		uc.State = 0
		err = ucr.SetState(ctx, uc.ID, 0)
	}
	if err != nil {
		fmt.Printf("start handler failed to save user-chat record in db: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "Failed, try again later",
		})
		if err != nil {
			fmt.Printf("start handler failed to send error-message: %v", err)
		}
	}

	// @TODO: Render Menu
}
