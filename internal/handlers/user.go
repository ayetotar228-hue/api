package handlers

import (
	"api/internal/broker"
	"api/internal/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]models.User, error)
	Create(ctx context.Context, email string) (*models.User, error)
}

type UserHandler struct {
	repo     UserRepository
	producer broker.ProducerInterface
}

func NewUserHandler(repo UserRepository, producer broker.ProducerInterface) *UserHandler {
	return &UserHandler{
		repo:     repo,
		producer: producer,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	user, err := h.repo.Create(ctx, req.Email)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.producer.SendUserCreated(ctx, user.ID, user.Email); err != nil {
		log.Printf("Failed to send broker event: %v", err)
	}

	log.Printf("User created: %d, event sent to broker", user.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
