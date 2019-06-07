package content

import (
	"github.com/yaches/habr_crawler/models"
)

type Storage interface {
	AddUsers([]models.User) error
	AddPosts([]models.Post) error
	AddComments([]models.Comment) error
}
