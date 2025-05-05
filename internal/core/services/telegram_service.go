package services

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/repositories"
	"math/rand"
	"sync"
	"time"
)

type TelegramService struct {
	editBot *bot.Bot
	readBot *bot.Bot
}

var (
	telegramServiceOnce     *sync.Once
	telegramServiceInstance *TelegramService
)

func InitializeTelegramService(editBot *bot.Bot, readBot *bot.Bot) {
	telegramServiceOnce.Do(func() {
		telegramServiceInstance = &TelegramService{
			editBot: editBot,
			readBot: readBot,
		}
	})
}

func GetTelegramService() *TelegramService {
	return telegramServiceInstance
}

func (ts *TelegramService) PublishPost(post *models.PostModel) {
	ur := repositories.GetUserRepository()

	msgTxt := fmt.Sprintf("%s\n\n%s", post.Title, post.Body)

	offset := 0
	var users = make([]*models.UserModel, 0)
	var err error = nil
	var duration = time.Duration(1) * time.Second
	for {
		users, err = ur.ListUsers(offset, 100)
		offset = offset + 100
		if err != nil {
			fmt.Printf("TelegramService.PublishPost: %v", err)
			return
		}

		if len(users) < 1 {
			break
		}

		for _, user := range users {
			_ = ts.sendToUser(user, msgTxt)

			duration = time.Duration(1+rand.Intn(29)) * time.Second
			time.Sleep(duration)
		}
	}
}

func (ts *TelegramService) sendToUser(
	user *models.UserModel,
	msgTxt string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := ts.readBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: user.ChatID,
		Text:   msgTxt,
	})
	if err != nil {
		fmt.Printf("TelegramService::sendToUser: failed to send message, error: %v", err)
		return err
	}

	return nil
}

func (ts *TelegramService) SendLastPostsToUser(
	user *models.UserModel,
) error {
	pr := repositories.GetPostRepository()
	posts, err := pr.ListLastPosts(0, 5, models.POST_STATUS_PUBLISHED)
	if err != nil {
		fmt.Printf("TelegramService::SendLastPostsToUser: %v", err)
	}

	var duration = time.Duration(1) * time.Second
	for _, post := range posts {
		msgTxt := fmt.Sprintf("%s\n\n%s", post.Title, post.Body)
		_ = ts.sendToUser(user, msgTxt)

		duration = time.Duration(1+rand.Intn(29)) * time.Second
		time.Sleep(duration)
	}

	return nil
}
