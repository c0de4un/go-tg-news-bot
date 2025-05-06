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
	cfg     *models.TelegramConfig
}

var (
	telegramServiceOnce     sync.Once
	telegramServiceInstance *TelegramService
)

func InitializeTelegramService(
	editBot *bot.Bot,
	readBot *bot.Bot,
	cfg *models.TelegramConfig,
) {
	telegramServiceOnce.Do(func() {
		telegramServiceInstance = &TelegramService{
			editBot: editBot,
			readBot: readBot,
			cfg:     cfg,
		}
	})
}

func GetTelegramService() *TelegramService {
	return telegramServiceInstance
}

func GetEditBotID() int64 {
	return telegramServiceInstance.cfg.EditorBotID
}

func GetReadBotID() int64 {
	return telegramServiceInstance.cfg.ReaderBotID
}

func (ts *TelegramService) IsAdmin(tgID int64) bool {
	return tgID == ts.cfg.AdminID
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

	ucr := repositories.GetUserChatRepository()
	uc, err := ucr.GetUserChat(user.ID, models.CHAT_TYPE_READER, GetReadBotID())
	if err != nil {
		fmt.Printf("TelegramService::sendToUser: failed to find user-chat, error: %v", err)
		return err
	}

	_, err = ts.readBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: uc.ChatID,
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

func (ts *TelegramService) NotifyAdmin(msgTxt string) {
	ur := repositories.GetUserRepository()
	admin, err := ur.GetUserByTelegramID(ts.cfg.AdminID)
	if err != nil {
		fmt.Printf("TelegramService::sendToUser: failed to retrieve admin, error: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	ucr := repositories.GetUserChatRepository()
	uc, err := ucr.GetUserChat(admin.ID, models.CHAT_TYPE_EDITOR, GetEditBotID())
	if err != nil {
		fmt.Printf("TelegramService::sendToUser: failed to retrieve admin chat, error: %v", err)
		return
	}

	_, err = ts.editBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: uc.ChatID,
		Text:   msgTxt,
	})
	if err != nil {
		fmt.Printf("TelegramService::sendToUser: failed to send admin message, error: %v", err)
	}
}
