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

type PostRepository struct {
}

var (
	postRepositoryOnce     sync.Once
	postRepositoryInstance *PostRepository
)

func InitializePostRepository() {
	postRepositoryOnce.Do(func() {
		postRepositoryInstance = &PostRepository{}
	})
}

func GetPostRepository() *PostRepository {
	return postRepositoryInstance
}

func (pr *PostRepository) GetPostByUserAndStatus(userID int64, status int64) (*models.PostModel, error) {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("PostRepository.GetPostByUserAndStatus: dbm is nil, maybe app is terminating")
		return nil, fmt.Errorf("PostRepository.GetPostByUserAndStatus: dbm is nil, maybe app is terminating")
	}

	db := dbm.GetDBConnection()
	query := `
        SELECT *
        FROM posts 
        WHERE created_by = $1
        AND status = $2
        ORDER BY id DESC`

	p := &models.PostModel{}
	err := db.QueryRow(query, userID, status).Scan(
		&p.ID,
		&p.Title,
		&p.Body,
		&p.Status,
		&p.UUID,
		&p.CreatedBy,
		&p.PublishedAt,
		&p.CreatedAt,
		&p.PublishedAt,
		&p.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		fmt.Printf("PostRepository.GetPostByUserAndStatus: %s", err)
		return nil, err
	}

	return p, nil
}

func (pr *PostRepository) SetTitle(ctx context.Context, postID int64, title string) error {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("PostRepository.SetTitle: dbm is nil, maybe app is terminating")
		return fmt.Errorf("PostRepository.SetTitle: dbm is nil, maybe app is terminating")
	}
	db := dbm.GetDBConnection()

	query := `
        UPDATE posts
        SET title = $1
        WHERE id = $2;`
	_, err := db.ExecContext(ctx, query,
		title,
		postID,
	)

	return err
}
