package handlers

import (
	"api/internal/broker"
	"api/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*models.User), args.Error(1)
}

type MockbrokerProducer struct {
	mock.Mock
}

func (m *MockbrokerProducer) SendUserCreated(ctx context.Context, userID int, email string) error {
	args := m.Called(ctx, userID, email)
	return args.Error(0)
}

func (m *MockbrokerProducer) SendEmailNotification(ctx context.Context, to, subject, body string) error {
	args := m.Called(ctx, to, subject, body)
	return args.Error(0)
}

func (m *MockbrokerProducer) SendMessage(ctx context.Context, topic string, msg broker.Message) error {
	args := m.Called(ctx, topic, msg)
	return args.Error(0)
}

func (m *MockbrokerProducer) Close() error {
	return nil
}

func TestUserHandler_GetUsers_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockProducer := new(MockbrokerProducer)

	expectedUsers := []models.User{
		{ID: 1, Email: "test1@example.com"},
		{ID: 2, Email: "test2@example.com"},
	}

	mockRepo.On("GetAll", mock.Anything).Return(expectedUsers, nil)

	handler := NewUserHandler(mockRepo, mockProducer)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	handler.GetUsers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var users []models.User
	err := json.NewDecoder(w.Body).Decode(&users)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "test1@example.com", users[0].Email)

	mockRepo.AssertExpectations(t)
}

func TestUserHandler_GetUsers_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockProducer := new(MockbrokerProducer)

	mockRepo.On("GetAll", mock.Anything).Return([]models.User{}, assert.AnError)

	handler := NewUserHandler(mockRepo, mockProducer)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	handler.GetUsers(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockProducer := new(MockbrokerProducer)

	newUser := &models.User{
		ID:    1,
		Email: "newuser@example.com",
	}

	mockRepo.On("Create", mock.Anything, "newuser@example.com").Return(newUser, nil)
	mockProducer.On("SendUserCreated", mock.Anything, 1, "newuser@example.com").Return(nil)

	handler := NewUserHandler(mockRepo, mockProducer)

	body := bytes.NewBufferString(`{"email":"newuser@example.com"}`)
	req := httptest.NewRequest("POST", "/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var user models.User
	err := json.NewDecoder(w.Body).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, "newuser@example.com", user.Email)

	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestUserHandler_CreateUser_InvalidJSON(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockProducer := new(MockbrokerProducer)

	handler := NewUserHandler(mockRepo, mockProducer)

	body := bytes.NewBufferString(`{invalid json}`)
	req := httptest.NewRequest("POST", "/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	mockProducer.AssertNotCalled(t, "SendUserCreated", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserHandler_CreateUser_EmptyEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockProducer := new(MockbrokerProducer)

	handler := NewUserHandler(mockRepo, mockProducer)

	body := bytes.NewBufferString(`{"email":""}`)
	req := httptest.NewRequest("POST", "/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	mockProducer.AssertNotCalled(t, "SendUserCreated", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserHandler_CreateUser_RepoError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockProducer := new(MockbrokerProducer)

	mockRepo.On("Create", mock.Anything, "test@example.com").Return((*models.User)(nil), assert.AnError)

	handler := NewUserHandler(mockRepo, mockProducer)

	body := bytes.NewBufferString(`{"email":"test@example.com"}`)
	req := httptest.NewRequest("POST", "/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockRepo.AssertExpectations(t)
	mockProducer.AssertNotCalled(t, "SendUserCreated", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserHandler_CreateUser_brokerError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockProducer := new(MockbrokerProducer)

	newUser := &models.User{
		ID:    1,
		Email: "test@example.com",
	}

	mockRepo.On("Create", mock.Anything, "test@example.com").Return(newUser, nil)
	mockProducer.On("SendUserCreated", mock.Anything, 1, "test@example.com").Return(assert.AnError)

	handler := NewUserHandler(mockRepo, mockProducer)

	body := bytes.NewBufferString(`{"email":"test@example.com"}`)
	req := httptest.NewRequest("POST", "/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var user models.User
	err := json.NewDecoder(w.Body).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)

	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}
