package services

import "github.com/samuelhgf/golang-study-restful-api/src/models"

type PostService interface {
	CreatePost(request *models.CreatePostRequest) (*models.DBPost, error)
	UpdatePost(string, *models.UpdatePost) (*models.DBPost, error)
	FindPostById(string) (*models.DBPost, error)
	FindPosts(page int, limit int) ([]*models.DBPost, error)
	DeletePost(string) error
}
