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

func (pr *PostRepository) GetPostByStatus(status int64) (*models.PostModel, error) {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("PostRepository.GetPostByStatus: dbm is nil, maybe app is terminating")
		return nil, fmt.Errorf("PostRepository.GetPostByStatus: dbm is nil, maybe app is terminating")
	}

	db := dbm.GetDBConnection()
	query := `
        SELECT *
        FROM posts 
        WHERE status = $2
        ORDER BY id DESC`

	p := &models.PostModel{}
	err := db.QueryRow(query, status).Scan(
		&p.ID,
		&p.Title,
		&p.Body,
		&p.Status,
		&p.CreatedBy,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		fmt.Printf("PostRepository.GetPostByStatus: %s", err)
		return nil, err
	}

	return p, nil
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
		&p.CreatedBy,
		&p.CreatedAt,
		&p.UpdatedAt,
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

func (pr *PostRepository) Create(ctx context.Context, userID int64) (*models.PostModel, error) {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("PostRepository.Create: dbm is nil, maybe app is terminating")
		return nil, fmt.Errorf("PostRepository.Create: dbm is nil, maybe app is terminating")
	}
	db := dbm.GetDBConnection()

	now := time.Now()

	query := `
        INSERT INTO posts (
            title,
            body,
            status,
            created_by,
            created_at,
            updated_at,
            deleted_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, title, body, status, created_by, created_at, updated_at, deleted_at`

	post := &models.PostModel{}
	err := db.QueryRowContext(
		ctx,
		query,
		"untitled",
		"empty",
		models.POST_STATUS_DRAFT,
		userID,
		now,
		now,
		now,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Body,
		&post.Status,
		&post.CreatedBy,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		fmt.Printf("PostRepository.Create: %s", err)
		return nil, err
	}

	return post, nil
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

func (pr *PostRepository) SetBodyWithState(ctx context.Context, postID int64, body string, status int64) error {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("PostRepository.SetTitle: dbm is nil, maybe app is terminating")
		return fmt.Errorf("PostRepository.SetTitle: dbm is nil, maybe app is terminating")
	}
	db := dbm.GetDBConnection()

	query := `
        UPDATE posts
        SET body = $1, status = $2
        WHERE id = $3;`
	_, err := db.ExecContext(ctx, query,
		body,
		status,
		postID,
	)

	return err
}

func (pr *PostRepository) ListLastPosts(offset, limit int, status int64) ([]*models.PostModel, error) {
	dbm, _ := database.GetDBManager()
	if dbm == nil {
		fmt.Println("PostRepository.ListPosts: dbm is nil, maybe app is terminating")
		return nil, fmt.Errorf("database manager not available")
	}
	db := dbm.GetDBConnection()

	query := `
        SELECT 
            id,
            title,
            body,
            status,
            created_by,
            created_at,
            updated_at,
            deleted_at
        FROM posts
        WHERE status = $1
        ORDER BY created_at DESC
        OFFSET $2
        LIMIT $3`

	rows, err := db.Query(query, status, offset, limit)
	if err != nil {
		fmt.Printf("PostRepository.ListPosts: query error: %v\n", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var posts []*models.PostModel

	for rows.Next() {
		post := &models.PostModel{}
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Body,
			&post.Status,
			&post.CreatedBy,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.DeletedAt,
		)
		if err != nil {
			fmt.Printf("PostRepository.ListPosts: scan error: %v\n", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("PostRepository.ListPosts: rows error: %v\n", err)
		return nil, err
	}

	return posts, nil
}
