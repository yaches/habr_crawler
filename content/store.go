package content

import "github.com/yaches/habr_crawler/models"

type Storage interface {
	AddUsers([]models.User) error
	GetAllUsers() ([]models.User, error)

	AddPosts([]models.Post) error
	GetAllPosts() ([]models.Post, error)

	AddComments([]models.Comment) error
	GetAllComments() ([]models.Comment, error)
}
