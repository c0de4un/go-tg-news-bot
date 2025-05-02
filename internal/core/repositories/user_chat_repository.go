package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	database "gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/databse"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
	"sync"
	"time"
)

type UserChatRepository struct {
}

var (
	userChatRepositoryOnce     sync.Once
	userChatRepositoryInstance *UserChatRepository
)

func InitializeUserChatRepository() {
	userChatRepositoryOnce.Do(func() {
		userChatRepositoryInstance = &UserChatRepository{}
	})
}

func GetUserChatRepository() *UserChatRepository {
	return userChatRepositoryInstance
}

func (ucr *UserChatRepository) GetUserChat(userID int64, chatID int64) (*models.ChatModel, error) {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("UserChatRepository.GetUserChat: dbm is nil, maybe app is terminating")
		return nil, fmt.Errorf("UserChatRepository.GetUserChat: dbm is nil, maybe app is terminating")
	}

	db := dbm.GetDBConnection()
	query := `
        SELECT *
        FROM user_chats 
        WHERE user_id = $1
        AND chat_id = $2`

	uc := &models.ChatModel{}
	err := db.QueryRow(query, userID, chatID).Scan(
		&uc.ID,
		&uc.UserID,
		&uc.ChatID,
		&uc.State,
		&uc.CreatedAt,
		&uc.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		fmt.Printf("UserChatRepository.GetUserChat: %s", err)
		return nil, err
	}

	return uc, nil
}

func (ucr *UserChatRepository) CreateUserChat(ctx context.Context, userID int64, chatID int64) (*models.ChatModel, error) {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("UserChatRepository.CreateUserChat: dbm is nil, maybe app is terminating")
		return nil, fmt.Errorf("UserChatRepository.CreateUserChat: dbm is nil, maybe app is terminating")
	}
	db := dbm.GetDBConnection()

	now := time.Now()

	query := `
        INSERT INTO user_chats (
            user_id,
            chat_id,
            state,
            created_at,
            updated_at
        ) VALUES ($1, $2, $3, $4, $5)
        RETURNING id, user_id, chat_id, state, created_at, updated_at`

	uc := &models.ChatModel{}
	err := db.QueryRowContext(ctx, query,
		userID,
		chatID,
		0,
		now,
		now,
	).Scan(
		&uc.ID,
		&uc.UserID,
		&uc.ChatID,
		&uc.State,
		&uc.CreatedAt,
		&uc.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		fmt.Printf("UserChatRepository.CreateUserChat: %s", err)
		return nil, err
	}

	return uc, nil
}

func (ucr *UserChatRepository) SetState(ctx context.Context, userChatID int64, state int64) error {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("UserChatRepository.CreateUserChat: dbm is nil, maybe app is terminating")
		return fmt.Errorf("UserChatRepository.CreateUserChat: dbm is nil, maybe app is terminating")
	}
	db := dbm.GetDBConnection()

	query := `
        UPDATE user_chats
        SET state = $1
        WHERE id = $2;`
	_, err := db.ExecContext(ctx, query,
		state,
		userChatID,
	)

	return err
}
