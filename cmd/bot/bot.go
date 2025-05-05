package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	database "gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/databse"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/repositories"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/services"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/tg"
	"os"
	"os/signal"
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
	repositories.InitializeUserChatRepository()
	repositories.InitializeClientRepository()
	repositories.InitializePostRepository()
	services.InitializeTranslateService()

	opts := []bot.Option{
		bot.WithDefaultHandler(tg.DefaultHandler),
		bot.WithDebug(),
	}

	editBot, err := bot.New(cfg.EditorToken, opts...)
	if err != nil {
		panic(err)
	}
	editBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)

	opts = []bot.Option{
		bot.WithDefaultHandler(tg.ReadBotDefaultHandler),
		bot.WithDebug(),
	}
	readBot, err := bot.New(cfg.ReaderToken, opts...)
	if err != nil {
		panic(err)
	}
	readBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, tg.ReadBotStartHandler)

	services.InitializeTelegramService(editBot, readBot, cfg.AdminID)

	go readBot.Start(ctx)

	editBot.Start(ctx)
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	tg.ClientStartHandler(ctx, b, update)
}
