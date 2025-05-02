package tg

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	newsmodels "gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
)

func handlePostCreateStart(user *newsmodels.UserModel, ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	// @TODO: handlePostCreateStart()
}
