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

func (ur *UserRepository) ListUsers(offset, limit int) ([]*models.UserModel, error) {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("UserRepository.ListUsers: dbm is nil, maybe app is terminating")
		return nil, fmt.Errorf("database manager not available")
	}
	db := dbm.GetDBConnection()

	query := `
        SELECT 
            users.id,
            users.telegram_username,
            users.telegram_id,
            users.uuid,
            users.created_at,
            users.updated_at,
            users.role,
            user_chats.chat_id
        FROM users
        JOIN user_chats ON user_chats.user_id = users.id
        ORDER BY users.created_at DESC
        OFFSET $1
        LIMIT $2`

	rows, err := db.Query(query, offset, limit)
	if err != nil {
		fmt.Printf("UserRepository.ListUsers: query error: %v\n", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var users []*models.UserModel

	for rows.Next() {
		user := &models.UserModel{}
		err := rows.Scan(
			&user.ID,
			&user.TelegramUsername,
			&user.TelegramID,
			&user.UUID,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Role,
			&user.ChatID,
		)
		if err != nil {
			fmt.Printf("UserRepository.ListUsers: scan error: %v\n", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("UserRepository.ListUsers: rows error: %v\n", err)
		return nil, err
	}

	return users, nil
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

func (ur *UserRepository) GetUserWithRelations(
	telegramID int64,
	chatType int64,
	chatID int64,
) (*models.UserModel, error) {
	user, err := ur.GetUserByTelegramID(telegramID)
	if err != nil {
		fmt.Printf("UserRepository.GetUserWithRelations: %v", err)
		return user, err
	}

	ucr := GetUserChatRepository()
	user.Chat, err = ucr.GetUserChat(user.ID, chatType, chatID)
	if err != nil {
		fmt.Printf("UserRepository.GetUserWithRelations: %v", err)
	}

	cr := GetClientRepository()
	user.Client, err = cr.GetClientByUserID(user.ID)
	if err != nil {
		fmt.Printf("UserRepository.GetUserWithRelations: %v", err)
	}

	return user, nil
}

func (ur *UserRepository) SetRole(ctx context.Context, userID int64, role int64) error {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("UserRepository.SetRole: dbm is nil, maybe app is terminating")
		return fmt.Errorf("database manager not available")
	}
	db := dbm.GetDBConnection()

	query := `
        UPDATE users
        SET role = $1,
            updated_at = $2
        WHERE id = $3`

	_, err := db.ExecContext(ctx, query,
		role,
		time.Now().UTC(), // Update timestamp
		userID,
	)

	if err != nil {
		fmt.Printf("UserRepository.SetRole: update error: %v\n", err)
		return fmt.Errorf("failed to update user role: %w", err)
	}

	return nil
}
