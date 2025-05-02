package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	database "gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/databse"
	"gitlab.com/korgi.tech/projects/go-news-tg-bot/internal/core/models"
	"sync"
)

type ClientRepository struct {
}

var (
	clientRepositoryOnce     sync.Once
	clientRepositoryInstance *ClientRepository
)

func InitializeClientRepository() {
	clientRepositoryOnce.Do(func() {
		clientRepositoryInstance = &ClientRepository{}
	})
}

func GetClientRepository() *ClientRepository {
	return clientRepositoryInstance
}

func (cr *ClientRepository) GetClientByUserID(userID int64) (*models.ClientModel, error) {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("ClientRepository.GetClientByUserID: dbm is nil, maybe app is terminating")
		return nil, fmt.Errorf("ClientRepository.GetClientByUserID: dbm is nil, maybe app is terminating")
	}

	db := dbm.GetDBConnection()
	query := `
        SELECT *
        FROM clients 
        WHERE user_id = $1`

	c := &models.ClientModel{}
	err := db.QueryRow(query, userID).Scan(
		&c.ID,
		&c.UserID,
		&c.Status,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		fmt.Printf("ClientRepository.GetClientByUserID: %s", err)
		return nil, err
	}

	return c, nil
}

func (cr *ClientRepository) Create(ctx context.Context, userID int64) (*models.ClientModel, error) {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("ClientRepository.Create: dbm is nil, maybe app is terminating")
		return nil, fmt.Errorf("ClientRepository.Create: dbm is nil, maybe app is terminating")
	}
	db := dbm.GetDBConnection()

	query := `
        INSERT INTO clients (
            user_id,
            status
        ) VALUES ($1, $2)
        RETURNING id, user_id, status, created_at, updated_at`

	c := &models.ClientModel{}
	err := db.QueryRowContext(ctx, query,
		userID,
		models.CLIENT_STATE_REGISTRATION,
	).Scan(
		&c.ID,
		&c.UserID,
		&c.Status,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		fmt.Printf("ClientRepository.Create: %s", err)
		return nil, err
	}

	return c, nil
}
