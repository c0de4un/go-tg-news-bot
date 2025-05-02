package repositories

import (
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
