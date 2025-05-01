package main

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	database "gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/databse"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/repositories"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/tg"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg, err := tg.LoadTelegramConfig()
	if err != nil {
		panic(err)
	}

	err = database.InitializeDB()
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	repositories.InitializeUserRepository()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithDebug(),
	}

	b, err := bot.New(cfg.Token, opts...)
	if err != nil {
		panic(err)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})

	if err != nil {
		fmt.Printf("default tg handler failed: %v", err)
	}
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	username := update.Message.From.Username
	chatId := update.Message.Chat.ID

	ur := repositories.GetUserRepository()
	user, err := ur.GetUserByTelegramID(userId)
	if err != nil {
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
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   fmt.Sprintf("Welcome, %s", username),
		})
		if err != nil {
			fmt.Printf("start handler failed: %v", err)
		}
	}
}
