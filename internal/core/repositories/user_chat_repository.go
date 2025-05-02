package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	database "gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/databse"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
	"sync"
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
