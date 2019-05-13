package contentstore

import "github.com/yaches/habr_crawler/models"

type ContentStore interface {
	AddUsers([]models.User) error
	GetAllUsers() ([]models.User, error)

	AddPosts([]models.Post) (int, error)
	GetAllPosts() ([]models.Post, error)

	AddComments([]models.Comment) (int, error)
	GetAllComments() ([]models.Comment, error)
}
