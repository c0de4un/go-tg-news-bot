package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"

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
		&user.Role,
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

func (ur *UserRepository) Create(ctx context.Context, telegramID int64, username string, role int64) (*models.UserModel, error) {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("UserRepository.Create: dbm is nil, maybe app is terminating")
		return nil, fmt.Errorf("UserRepository.Create: dbm is nil, maybe app is terminating")
	}
	db := dbm.GetDBConnection()

	newUUID := uuid.New().String()
	now := time.Now()

	query := `
        INSERT INTO users (
            telegram_username,
            telegram_id,
            uuid,
            created_at,
            updated_at,
            role
        ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, telegram_username, telegram_id, uuid, created_at, updated_at, role`

	user := &models.UserModel{}
	err := db.QueryRowContext(ctx, query,
		username,
		telegramID,
		newUUID,
		now,
		now,
		role,
	).Scan(
		&user.ID,
		&user.TelegramUsername,
		&user.TelegramID,
		&user.UUID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		fmt.Printf("UserRepository.Create: %s", err)
		return nil, err
	}

	return user, nil
}
