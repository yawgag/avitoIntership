package authService

import (
	"context"
	"errors"
	"orderPickupPoint/config"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/storage"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo      storage.Auth
	cfg           *config.Config
	tokensHandler AuthTokenHandler
}

type AuthTokenHandler interface {
	CreateAccessToken(ctx context.Context, user *models.User) (string, error)
	CreateRefreshToken(ctx context.Context, user *models.User) (string, error)
	ParseJwt(token string) (*jwt.MapClaims, error)
}
type TokenHandlerImpl struct {
	secretWord string
	authRepo   storage.Auth
}

func NewAuthService(authRepo storage.Auth, cfg *config.Config) *AuthService {
	handler := &TokenHandlerImpl{secretWord: cfg.SecretWord, authRepo: authRepo}
	return &AuthService{
		authRepo:      authRepo,
		cfg:           cfg,
		tokensHandler: handler,
	}
}

func (s *TokenHandlerImpl) CreateRefreshToken(ctx context.Context, user *models.User) (string, error) {
	sessionId := uuid.New().String()
	refreshTokenExpireTime, err := s.authRepo.CreateSession(ctx, user, sessionId)
	refreshTokenExpireTime = time.Now().Add(30 * 24 * time.Hour)
	if err != nil {
		return "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sessionId":  sessionId,
		"expireTime": refreshTokenExpireTime,
	})

	signedRefreshToken, err := refreshToken.SignedString([]byte(s.secretWord))
	if err != nil {
		return "", err
	}

	return signedRefreshToken, nil
}

func (s *TokenHandlerImpl) CreateAccessToken(ctx context.Context, user *models.User) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":     user.Id,
		"userRole":   user.Role,
		"expireTime": time.Now().Add(15 * time.Minute),
	})

	signedAccessToken, err := accessToken.SignedString([]byte(s.secretWord))
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}

func (s *AuthService) DummyLogin(ctx context.Context, user *models.User) (*models.AuthTokens, error) {
	refreshToken, err := s.tokensHandler.CreateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokensHandler.CreateAccessToken(ctx, user)
	if err != nil {
		return nil, err
	}

	return &models.AuthTokens{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, user *models.User) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(passwordHash)

	err = s.authRepo.AddNewUser(ctx, user)
	return err
}

func (s *AuthService) Login(ctx context.Context, user *models.User) (*models.AuthTokens, error) {
	userFromDb, err := s.authRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFromDb.PasswordHash), []byte(user.Password))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokensHandler.CreateRefreshToken(ctx, userFromDb)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokensHandler.CreateAccessToken(ctx, userFromDb)
	if err != nil {
		return nil, err
	}

	return &models.AuthTokens{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

// verify the tokens and update them if necessary
func (s *AuthService) HandleTokens(ctx context.Context, tokens *models.AuthTokens) (*models.AuthTokens, error) {
	// take info from both tokens
	accessTokenClaims, err := s.tokensHandler.ParseJwt(tokens.AccessToken)
	if err != nil {
		return nil, err
	}
	refreshTokenClaims, err := s.tokensHandler.ParseJwt(tokens.RefreshToken)
	if err != nil {
		return nil, err
	}

	// access token expire time
	ATexpireTime, err := time.Parse(time.RFC3339, (*accessTokenClaims)["expireTime"].(string))
	if err != nil {
		return nil, err
	}

	NewTokens := &models.AuthTokens{}

	if ATexpireTime.Before(time.Now()) { // if access token is expired
		NewTokens.NewAccessToken = true
		RTexpireTime, err := time.Parse(time.RFC3339, (*refreshTokenClaims)["expireTime"].(string))
		if err != nil {
			return nil, err
		}

		if RTexpireTime.Before(time.Now()) {
			NewTokens.NewRefreshToken = true

			// update expire time of refresh token
			newExpireTime, err := s.authRepo.UpdateSessionExpireTime(ctx, (*refreshTokenClaims)["sessionId"].(string))
			newExpireTime = time.Now().Add(30 * 24 * time.Hour)
			if err != nil {
				return nil, err
			}

			// update time in old token struct and generate new token from this struct
			(*refreshTokenClaims)["expireTime"] = newExpireTime
			newRefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(s.cfg.SecretWord))
			if err != nil {
				return nil, err
			}
			NewTokens.RefreshToken = newRefreshToken

		}

		session, err := s.authRepo.GetSession(ctx, (*refreshTokenClaims)["sessionId"].(string))
		if err != nil {
			return nil, err
		}

		(*accessTokenClaims)["userRole"] = session.UserRole
		(*accessTokenClaims)["expireTime"] = time.Now().Add(time.Hour * 30 * 24)

		newAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString([]byte(s.cfg.SecretWord))
		if err != nil {
			return nil, err
		}

		NewTokens.AccessToken = newAccessToken
	}

	return NewTokens, nil
}
func (s *AuthService) AvaliableForUser(tokens *models.AuthTokens, avaliableRoles []string) (bool, error) {
	accessTokenClaims, err := s.tokensHandler.ParseJwt(tokens.AccessToken)
	if err != nil {
		return false, err
	}
	userRole := (*accessTokenClaims)["userRole"].(string)
	if slices.Contains(avaliableRoles, userRole) {
		return true, nil
	}
	return false, nil

}

func (p *TokenHandlerImpl) ParseJwt(token string) (*jwt.MapClaims, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("wrong token format")
		}
		return []byte(p.secretWord), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok && jwtToken.Valid {
		return &claims, nil
	}

	return nil, errors.New("wrong token")
}
