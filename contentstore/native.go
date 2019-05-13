package contentstore

import "github.com/yaches/habr_crawler/models"

type ContentStoreNative struct {
	users    []models.User
	posts    []models.Post
	comments []models.Comment
}

func NewContentStoreNative() ContentStoreNative {
	return ContentStoreNative{}
}

func (cs *ContentStoreNative) AddUsers(users []models.User) error {
	cs.users = append(cs.users, users...)
	return nil
}

func (cs *ContentStoreNative) GetAllUsers() ([]models.User, error) {
	return cs.users, nil
}

func (cs *ContentStoreNative) AddPosts(posts []models.Post) error {
	cs.posts = append(cs.posts, posts...)
	return nil
}

func (cs *ContentStoreNative) GetAllPosts() ([]models.Post, error) {
	return cs.posts, nil
}

func (cs *ContentStoreNative) AddComments(comments []models.Comment) error {
	cs.comments = append(cs.comments, comments...)
	return nil
}

func (cs *ContentStoreNative) GetAllComments() ([]models.Comment, error) {
	return cs.comments, nil
}
