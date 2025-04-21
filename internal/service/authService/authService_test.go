package authService

// import (
// 	"context"
// 	"orderPickupPoint/config"
// 	"orderPickupPoint/internal/models"
// 	"orderPickupPoint/internal/storage"
// 	"testing"

// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// 	"golang.org/x/crypto/bcrypt"
// )

// type MockAuthRepo struct {
// 	mock.Mock
// 	storage.Auth
// }

// func (m *MockAuthRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
// 	args := m.Called(ctx, email)
// 	return args.Get(0).(*models.User), args.Error(1)
// }

// type MockTokenHandler struct {
// 	mock.Mock
// }

// func (m *MockTokenHandler) CreateAccessToken(ctx context.Context, user *models.User) (string, error) {
// 	args := m.Called(ctx, user)
// 	return args.String(0), args.Error(1)
// }

// func (m *MockTokenHandler) CreateRefreshToken(ctx context.Context, user *models.User) (string, error) {
// 	args := m.Called(ctx, user)
// 	return args.String(0), args.Error(1)
// }

// func (m *MockTokenHandler) ParseJwt(token string) (*jwt.MapClaims, error) {
// 	args := m.Called(token)
// 	return args.Get(0).(*jwt.MapClaims), args.Error(1)
// }
