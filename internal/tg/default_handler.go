package tg

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	models2 "gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/repositories"
)

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	ur := repositories.GetUserRepository()
	ucr := repositories.GetUserChatRepository()
	user, err := ur.GetUserWithRelations(update.Message.From.ID)
	if err != nil {
		fmt.Printf("DefaultHandler: failed to retrieve user with error %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Failed, try again later",
		})
		if err != nil {
			fmt.Printf("DefaultHandler: failed to send error message with %v", err)
		}

		return
	}

	// Cancel, if not Client
	if user.Role != models2.USER_ROLE_CLIENT {
		fmt.Println("DefaultHandler: not a client, skip")
		return
	}

	pr := repositories.GetPostRepository()
	post, err := pr.GetPostByUserAndStatus(user.ID, models2.POST_STATUS_DRAFT)
	if err != nil {
		fmt.Printf("DefaultHandler: failed to retrieve post, with error: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Failed, try again later",
		})
		if err != nil {
			fmt.Printf("DefaultHandler: failed to send error message with %v", err)
		}

		return
	}

	// Title prompt
	if user.Chat.State == models2.CHAT_STATE_POST_TITLE {
		if len(update.Message.Text) < 3 {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Invalid title",
			})
			if err != nil {
				fmt.Printf("DefaultHandler: failed to send error message with %v", err)
			}

			return
		}

		err = pr.SetTitle(ctx, post.ID, update.Message.Text)
		if err != nil {
			fmt.Printf("DefaultHandler: failed to set chat state, with error: %v", err)
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Failed, try again later",
			})
			if err != nil {
				fmt.Printf("DefaultHandler: failed to send error message with %v", err)
			}

			return
		}

		err = ucr.SetState(ctx, user.Chat.ID, models2.CHAT_STATE_POST_BODY)
		if err != nil {
			fmt.Printf("DefaultHandler: failed to set chat state, with error: %v", err)
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Failed, try again later",
			})
			if err != nil {
				fmt.Printf("DefaultHandler: failed to send error message with %v", err)
			}

			return
		}

		return
	}

	if user.Chat.State == models2.CHAT_STATE_POST_BODY {
		// @TODO: Validate Body
		// @TODO: Reset Chat-State
		// @TODO: Change Post state to Published
	}
}
