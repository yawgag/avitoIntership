package service

import (
	"context"
	"orderPickupPoint/config"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/service/authService"
	"orderPickupPoint/internal/service/pickupPointService"
	"orderPickupPoint/internal/service/receptionService"
	"orderPickupPoint/internal/storage"

	"github.com/google/uuid"
)

type PickupPoint interface {
	// Create(ctx context.Context, data *models.PickupPoint) error
	// ??? GetAllInfo
	Create(ctx context.Context, pickupPoint *models.PickupPointAPI) (*models.PickupPointAPI, error)
}

type Reception interface {
	// Create(ctx context.Context, ... ) error
	// AddProduct(ctx context.Context, ...) error
	// DeleteProduct(ctx context.Context, ...) error
	// Close(ctx context.Context, ...) error
	CreateReception(ctx context.Context, pvzId uuid.UUID) (*models.ReceptionAPI, error)
	AddProduct(ctx context.Context, productAPI *models.ProductAPI) (*models.ProductAPI, error)
}

type Auth interface {
	CreateAccessToken(ctx context.Context, user *models.User) (string, error)
	CreateRefreshToken(ctx context.Context, user *models.User) (string, error)
	Register(ctx context.Context, user *models.User) error
	Login(ctx context.Context, user *models.User) (*models.AuthTokens, error)

	AvaliableForUser(tokens *models.AuthTokens, avaliableRoles []string) (bool, error)
	HandleTokens(ctx context.Context, tokens *models.AuthTokens) (*models.AuthTokens, error)
}

type Deps struct {
	Repos *storage.Repositories
	Cfg   *config.Config
}

type Services struct {
	PickupPoint PickupPoint
	Reception   Reception
	Auth        Auth
}

func NewServices(deps *Deps) *Services {
	return &Services{
		PickupPoint: pickupPointService.NewPickupPointService(deps.Repos.PickupPoint),
		Reception:   receptionService.NewReceptionService(deps.Repos.Reception),
		Auth:        authService.NewAuthService(deps.Repos.Auth, deps.Cfg),
	}
}
