package tg

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	models2 "gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/repositories"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/services"
)

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	ur := repositories.GetUserRepository()
	ucr := repositories.GetUserChatRepository()
	user, err := ur.GetUserWithRelations(update.Message.From.ID)
	if err != nil {
		fmt.Printf("DefaultHandler: failed to retrieve user with error %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   services.Translate("Failed, try again later"),
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
			Text:   services.Translate("Failed, try again later"),
		})
		if err != nil {
			fmt.Printf("DefaultHandler: failed to send error message with %v", err)
		}

		return
	}

	// Title prompt
	if user.Chat.State == models2.CHAT_STATE_POST_TITLE {
		ln := len(update.Message.Text)
		if ln < 3 || ln > 254 {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   services.Translate("Invalid title"),
			})
			if err != nil {
				fmt.Printf("DefaultHandler: failed to send error message with %v", err)
			}

			return
		}

		post, err := pr.Create(ctx, user.ID)
		if err != nil {
			fmt.Printf("DefaultHandler: failed to create new post, with error: %v", err)
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   services.Translate("Failed, try again later"),
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
				Text:   services.Translate("Failed, try again later"),
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
				Text:   services.Translate("Failed, try again later"),
			})
			if err != nil {
				fmt.Printf("DefaultHandler: failed to send error message with %v", err)
			}

			return
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   services.Translate("Enter main text"),
		})
		if err != nil {
			fmt.Printf("DefaultHandler: failed to send error message with %v", err)
		}

		return
	}

	if user.Chat.State == models2.CHAT_STATE_POST_BODY {
		ln := len(update.Message.Text)
		if ln < 3 || ln > 10000 {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   services.Translate("Invalid content"),
			})
			if err != nil {
				fmt.Printf("DefaultHandler: failed to send error message with %v", err)
			}

			return
		}

		err = pr.SetBodyWithState(ctx, post.ID, update.Message.Text, models2.POST_STATUS_PUBLISHED)
		if err != nil {
			fmt.Printf("DefaultHandler: failed to set post body, with error: %v", err)
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   services.Translate("Failed, try again later"),
			})
			if err != nil {
				fmt.Printf("DefaultHandler: failed to send error message with %v", err)
			}

			return
		}

		err = ucr.SetState(ctx, user.Chat.ID, models2.CHAT_STATE_POST_WELCOME)
		if err != nil {
			fmt.Printf("DefaultHandler: failed to set chat state, with error: %v", err)
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   services.Translate("Failed, try again later"),
			})
			if err != nil {
				fmt.Printf("DefaultHandler: failed to send error message with %v", err)
			}

			return
		}

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   services.Translate("Published"),
		})
	}
}
