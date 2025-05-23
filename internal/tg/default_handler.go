package tg

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	models2 "gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/repositories"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/services"
)

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	ur := repositories.GetUserRepository()
	if ur == nil {
		panic("GetUserRepository returned nil")
	}
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

	// Forward-Post Input
	if user.Chat.State == models2.CHAT_STATE_FORWARDED_POST_INPUT {
		fr := repositories.GetForwardPostRepository()
		fwdPost, err := fr.Create(user.ID, int64(update.Message.ID), update.Message.Chat.ID)
		if err != nil {
			fmt.Printf("DefaultHandler: failed to create forwarded message with %v", err)
		}

		ts := services.GetTelegramService()
		ts.NotifyAdminAboutNewPost(fwdPost)

		kb := inline.New(b).
			Row().
			Button(services.Translate("Create post"), []byte(fmt.Sprintf("1-1;%d", user.ID)), onInlineKeyboardSelect)
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        services.Translate("Sent to moderation"),
			ReplyMarkup: kb,
		})
		if err != nil {
			fmt.Printf("DefaultHandler: failed to send message with %v", err)
		}

		return
	}
}
