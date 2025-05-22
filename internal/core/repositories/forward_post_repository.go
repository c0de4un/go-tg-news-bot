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
	fromChatID int64,
) (*models.ForwardPostModel, error) {
	now := time.Now()
	post := &models.ForwardPostModel{
		TelegramID: tgMsgID,
		FromChatID: fromChatID,
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
		   from_chat_id,
		   status,
		   created_by,
           created_at,
           updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := db.QueryRowContext(
		ctx,
		query,
		post.TelegramID,
		post.FromChatID,
		post.Status,
		post.CreatedBy,
		now,
		now,
	).Scan(
		&post.ID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		fmt.Printf("ForwardPostRepository.Create: %s", err)
		return nil, err
	}

	return post, nil
}

func (r *ForwardPostRepository) GetById(id int64) (*models.ForwardPostModel, error) {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("ForwardPostRepository.GetById: dbm is nil, maybe app is terminating")
		return nil, fmt.Errorf("ForwardPostRepository.GetById: dbm is nil, maybe app is terminating")
	}

	db := dbm.GetDBConnection()
	query := `
        SELECT *
        FROM forward_posts 
        WHERE id = $1
        LIMIT 1`

	row := db.QueryRow(query, id)

	return r.fillModelByRow(row, nil)
}

func (r *ForwardPostRepository) SetState(id int64, status int64) error {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("\nForwardPostRepository.SetState: dbm is nil, maybe app is terminating")
		return fmt.Errorf("ForwardPostRepository.SetState: dbm is nil, maybe app is terminating")
	}
	db := dbm.GetDBConnection()

	query := `
        UPDATE forward_posts
        SET status = $1
        WHERE id = $2;`
	_, err := db.Exec(query,
		status,
		id,
	)

	return err
}

func (r *ForwardPostRepository) GetByStatus(status int64) (*models.ForwardPostModel, error) {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("ForwardPostRepository.GetByStatus: dbm is nil, maybe app is terminating")
		return nil, fmt.Errorf("ForwardPostRepository.GetByStatus: dbm is nil, maybe app is terminating")
	}

	db := dbm.GetDBConnection()
	query := `
        SELECT *
        FROM forward_posts 
        WHERE status = $1
        LIMIT 1`

	row := db.QueryRow(query, status)

	return r.fillModelByRow(row, nil)
}

func (r *ForwardPostRepository) fillModelByRow(row *sql.Row, post *models.ForwardPostModel) (*models.ForwardPostModel, error) {
	if post == nil {
		post = &models.ForwardPostModel{}
	}

	err := row.Scan(
		&post.ID,
		&post.TelegramID,
		&post.FromChatID,
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
