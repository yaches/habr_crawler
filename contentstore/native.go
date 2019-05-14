package contentstore

import "github.com/yaches/habr_crawler/models"

type StorageNative struct {
	users    []models.User
	posts    []models.Post
	comments []models.Comment
}

func NewStorageNative() StorageNative {
	return StorageNative{}
}

func (cs *StorageNative) AddUsers(users []models.User) error {
	cs.users = append(cs.users, users...)
	return nil
}

func (cs *StorageNative) GetAllUsers() ([]models.User, error) {
	return cs.users, nil
}

func (cs *StorageNative) AddPosts(posts []models.Post) error {
	cs.posts = append(cs.posts, posts...)
	return nil
}

func (cs *StorageNative) GetAllPosts() ([]models.Post, error) {
	return cs.posts, nil
}

func (cs *StorageNative) AddComments(comments []models.Comment) error {
	cs.comments = append(cs.comments, comments...)
	return nil
}

func (cs *StorageNative) GetAllComments() ([]models.Comment, error) {
	return cs.comments, nil
}
