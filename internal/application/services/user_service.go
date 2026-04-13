package services

import (
	"api/internal/application/models"
	"api/internal/domain/user"
	"api/pkg/broker/producer"
	"context"
	"log"
)

type UserService struct {
	repo     user.Repository
	producer producer.Producer
}

func NewUserService(repo user.Repository, producer producer.Producer) *UserService {
	return &UserService{
		repo:     repo,
		producer: producer,
	}
}

func (s *UserService) CreateUser(ctx context.Context, email string) (*models.UserResponse, error) {
	u, err := user.NewUser(email)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	if err := s.producer.SendUserCreated(ctx, u.ID, u.Email); err != nil {
		log.Printf("Producer error: %v", err)
	}

	response := models.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}

	return &response, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.UserResponse, error) {

	users, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	response := make([]models.UserResponse, len(users))
	for i, u := range users {
		response[i] = models.UserResponse{
			ID:        u.ID,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		}
	}

	return response, nil
}
