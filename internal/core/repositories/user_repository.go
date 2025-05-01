package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/databse"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
)

type UserRepository struct {
}

var (
	userRepositoryOnce     sync.Once
	userRepositoryInstance *UserRepository
)

func InitializeUserRepository() {
	userRepositoryOnce.Do(func() {
		userRepositoryInstance = &UserRepository{}
	})
}

func GetUserRepository() *UserRepository {
	return userRepositoryInstance
}

func (ur *UserRepository) GetUserByTelegramID(telegramID int64) (*models.UserModel, error) {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("UserRepository.GetUserByTelegramID: dbm is nil, maybe app is terminating")
		return nil, fmt.Errorf("UserRepository.GetUserByTelegramID: dbm is nil, maybe app is terminating")
	}

	db := dbm.GetDBConnection()
	query := `
        SELECT *
        FROM users 
        WHERE telegram_id = $1`

	user := &models.UserModel{}
	err := db.QueryRow(query, telegramID).Scan(
		&user.ID,
		&user.TelegramUsername,
		&user.TelegramID,
		&user.UUID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		fmt.Printf("UserRepository.GetUserByTelegramID: %s", err)
		return nil, err
	}

	return user, nil
}
