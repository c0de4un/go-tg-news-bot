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
	user, err := ur.GetUserWithRelations(update.Message.From.ID, models2.CHAT_TYPE_EDITOR, services.GetEditBotID())
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
	if user.Role != models2.USER_ROLE_CLIENT && user.Role != models2.USER_ROLE_ADMIN {
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

	// Forward-Post Input
	if user.Chat.State == models2.CHAT_STATE_FORWARDED_POST_INPUT {
		ln := len(update.Message.Text)
		if ln < 3 || ln > 254 {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   services.Translate("Invalid content"),
			})
			if err != nil {
				fmt.Printf("DefaultHandler: failed to send error message with %v", err)
			}

			return
		}

		// @TODO: Store message and chat id as Forward-Post
		// @TODO: #1 Get ForwardPostRepository
		// @TODO: #2 Create ForwardPost
		ts := services.GetTelegramService()
		ts.NotifyAdminAboutNewPost(
			user,
			fwdPost,
		)

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   services.Translate("Sent to moderation"),
		})
		if err != nil {
			fmt.Printf("DefaultHandler: failed to send message with %v", err)
		}

		return
	}
}
