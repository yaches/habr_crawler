package content

import (
	"bytes"
	"encoding/json"

	"github.com/elastic/go-elasticsearch"
	"github.com/yaches/habr_crawler/models"
)

const (
	userIndex    = "users"
	postIndex    = "posts"
	commentIndex = "comments"
)

type StorageES struct {
	es *elasticsearch.Client
}

func NewStorageES() (*StorageES, error) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}

	return &StorageES{es: es}, nil
}

func (s *StorageES) AddUsers(users []models.User) error {
	for _, u := range users {
		if err := s.insertOne(userIndex, u.Username, u); err != nil {
			return err
		}
	}
	return nil
}

func (s *StorageES) GetAllUsers() ([]models.User, error) {
	return nil, nil
}

func (s *StorageES) AddPosts(posts []models.Post) error {
	for _, p := range posts {
		if err := s.insertOne(postIndex, p.ID, p); err != nil {
			return err
		}
	}
	return nil
}

func (s *StorageES) GetAllPosts() ([]models.Post, error) {
	return nil, nil
}

func (s *StorageES) AddComments(comments []models.Comment) error {
	for _, c := range comments {
		if err := s.insertOne(commentIndex, c.ID, c); err != nil {
			return err
		}
	}
	return nil
}

func (s *StorageES) GetAllComments() ([]models.Comment, error) {
	return nil, nil
}

func (s *StorageES) insertOne(index, id string, obj interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(obj); err != nil {
		return err
	}
	res, err := s.es.Create(index, id, &buf)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
