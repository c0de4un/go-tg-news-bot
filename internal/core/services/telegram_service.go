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

func (ts *TelegramService) PublishForwardPost(post *models.ForwardPostModel) {
	ur := repositories.GetUserRepository()
	ucr := repositories.GetUserChatRepository()
	ctx := context.Background()

	offset := 0
	var users = make([]*models.UserModel, 0)
	var err error = nil
	var duration = time.Duration(1) * time.Second
	sentUsers := make(map[int64]bool)
	for {
		users, err = ur.ListUsers(offset, 100)
		offset = offset + 100
		if err != nil {
			fmt.Printf("TelegramService.PublishForwardPost: %v", err)
			return
		}

		if len(users) < 1 {
			break
		}

		for _, user := range users {
			if sentUsers[user.ID] {
				continue
			}

			uc, err := ucr.GetUserChat(user.ID, models.CHAT_TYPE_READER, GetReadBotID())
			if err != nil || uc == nil {
				fmt.Printf("\nTelegramService::PublishForwardPost: failed to find user-chat, error: %v\n", err)
				continue
			}

			_, err = ts.editBot.CopyMessage(ctx, &bot.CopyMessageParams{
				ChatID:     ts.cfg.ChannelID,
				FromChatID: post.FromChatID,
				MessageID:  int(post.TelegramID),
			})

			if err != nil {
				fmt.Printf("\nTelegramService::PublishForwardPost: failed to send post for a user %d, error: %v\n", user.ID, err)
				continue
			}

			sentUsers[user.ID] = true

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

func (ts *TelegramService) SendForwardPost(
	post *models.ForwardPostModel,
	chatID int64,
	isReader bool,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var b *bot.Bot = nil
	if isReader {
		b = ts.readBot
	} else {
		b = ts.editBot
	}

	_, err := b.ForwardMessage(ctx, &bot.ForwardMessageParams{
		ChatID:     chatID,
		FromChatID: post.FromChatID,
		MessageID:  int(post.TelegramID),
	})
	if err != nil {
		fmt.Printf("\nTelegramService::SendForwardPost: failed to forward message, error: %v\n", err)
	}

	return err
}

func (ts *TelegramService) SendPost(
	ctx context.Context,
	post *models.PostModel,
	chatID int64,
	isReader bool,
) error {
	msgTxt := ts.renderPost(post)
	var b *bot.Bot = nil
	if isReader {
		b = ts.readBot
	} else {
		b = ts.editBot
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msgTxt,
	})

	return err
}

func (ts *TelegramService) NotifyAdminAboutNewPost(
	fwdPost *models.ForwardPostModel,
) {
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

	_, err = ts.editBot.ForwardMessage(ctx, &bot.ForwardMessageParams{
		ChatID:     uc.ChatID,
		FromChatID: uc.ChatID,
		MessageID:  int(fwdPost.TelegramID),
	})
	if err != nil {
		fmt.Printf("\nTelegramService::sendToUser: failed to forward message to admin, error: %v\n", err)
	}
}

func (ts *TelegramService) renderPost(post *models.PostModel) string {
	return fmt.Sprintf("%s\n\n%s", post.Title, post.Body)
}
