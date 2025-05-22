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

type ForwardPostRepository struct {
}

var (
	forwardPostRepositoryOnce     sync.Once
	forwardPostRepositoryInstance *ForwardPostRepository
)

func InitForwardPostRepository() {
	forwardPostRepositoryOnce.Do(func() {
		forwardPostRepositoryInstance = &ForwardPostRepository{}
	})
}

func GetForwardPostRepository() *ForwardPostRepository {
	return forwardPostRepositoryInstance
}

func (r *ForwardPostRepository) Create(
	userId int64,
	tgMsgID int64,
) (*models.ForwardPostModel, error) {
	now := time.Now()
	post := &models.ForwardPostModel{
		TelegramID: tgMsgID,
		Status:     models.POST_STATUS_MODERATION,
		CreatedBy:  userId,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("ForwardPostRepository.Create: dbm is nil, maybe app is terminating")
		return nil, fmt.Errorf("ForwardPostRepository.Create: dbm is nil, maybe app is terminating")
	}
	db := dbm.GetDBConnection()

	query := `
        INSERT INTO forward_posts (
		   telegram_id,
		   status,
		   created_by,
           created_at,
           updated_at
        ) VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	row := db.QueryRowContext(
		ctx,
		query,
		post.TelegramID,
		post.Status,
		post.CreatedBy,
		now,
		now,
	)

	return r.fillModelByRow(row, post)
}

func (r *ForwardPostRepository) fillModelByRow(row *sql.Row, post *models.ForwardPostModel) (*models.ForwardPostModel, error) {
	if post == nil {
		post = &models.ForwardPostModel{}
	}

	err := row.Scan(
		&post.ID,
		&post.TelegramID,
		&post.Status,
		&post.CreatedBy,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		fmt.Printf("ForwardPostRepository.fillModelByRow: %s", err)
		return nil, err
	}

	return post, nil
}
