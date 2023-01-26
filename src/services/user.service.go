package services

import "github.com/samuelhgf/golang-study-restful-api/src/models"

type UserService interface {
	FindUserById(string) (*models.DBResponse, error)
	FindUserByEmail(string) (*models.DBResponse, error)
}
