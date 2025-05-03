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
			Text:   services.Translate("Error, try again later"),
		})
		if err != nil {
			fmt.Printf("start handler failed: %v", err)
		}
	}
	if user != nil {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   fmt.Sprintf(services.Translate("Welcome back, %s"), user.TelegramUsername),
		})
		if err != nil {
			fmt.Printf("start handler failed: %v", err)
		}
	} else {
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

	cr := repositories.GetClientRepository()
	c, err := cr.GetClientByUserID(user.ID)
	if err != nil {
		fmt.Printf("start handler failed to retrieve client from db: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   services.Translate("Error, try again later"),
		})
		if err != nil {
			fmt.Printf("start handler failed to send error-message: %v", err)
		}
	}

	if c == nil {
		c, err = cr.Create(ctx, user.ID)
	}
	if err != nil {
		fmt.Printf("start handler failed to create client in db: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   services.Translate("Error, try again later"),
		})
		if err != nil {
			fmt.Printf("start handler failed to send error-message: %v", err)
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
	} else {
		uc.State = newsmodels.CHAT_STATE_POST_WELCOME
		err = ucr.SetState(ctx, uc.ID, newsmodels.CHAT_STATE_POST_WELCOME)
	}
	if err != nil {
		fmt.Printf("start handler failed to save user-chat record in db: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   services.Translate("Error, try again later"),
		})
		if err != nil {
			fmt.Printf("start handler failed to send error-message: %v", err)
		}
	}

	kb := inline.New(b).
		Row().
		Button(services.Translate("Create post"), []byte(fmt.Sprintf("1-1;%d", userId)), onInlineKeyboardSelect)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        services.Translate("Menu"),
		ReplyMarkup: kb,
	})
	if err != nil {
		fmt.Printf("start handler failed to save user-chat record in db: %v", err)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   services.Translate("Error, try again later"),
		})
		if err != nil {
			fmt.Printf("start handler failed to send error-message: %v", err)
		}
	}
}
