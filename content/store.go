package content

import (
	"github.com/yaches/habr_crawler/models"
	"github.com/yaches/habr_crawler/models/aggs"
)

type Storage interface {
	AddUsers([]models.User) error
	GetAllUsers() ([]models.User, error)

	AddPosts([]models.Post) error
	GetAllPosts() ([]models.Post, error)

	AddComments([]models.Comment) error
	GetAllComments() ([]models.Comment, error)

	GetCommonInfo() (aggs.CommonInfo, error)
}
