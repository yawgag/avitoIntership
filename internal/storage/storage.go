package storage

import (
	"context"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/storage/postgres"
	"orderPickupPoint/internal/storage/postgres/authRepo"
	"orderPickupPoint/internal/storage/postgres/pickupPointRepo"
	"orderPickupPoint/internal/storage/postgres/receptionRepo"
	"time"

	"github.com/google/uuid"
)

type PickupPoint interface {
	Create(ctx context.Context, pickupPoint *models.PickupPoint) (*models.PickupPoint, error)
	GetCityIdByName(ctx context.Context, name string) (int, error)
	GetFilteredInfo(ctx context.Context, filter *models.PvzFilter) ([]models.PvzFilteredInfo, error)
}

type Reception interface {
	CreateReception(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error)
	GetStatusNameById(ctx context.Context, id int) (string, error)
	GetProductTypeIdByName(ctx context.Context, name string) (int, error)
	AddProductToReception(ctx context.Context, product *models.Product, pvzId uuid.UUID) (*models.Product, error)
	DeleteLastProductInReception(ctx context.Context, pvzId uuid.UUID) error
	CloseReception(ctx context.Context, pvzId uuid.UUID) error
}

type Auth interface {
	CreateSession(ctx context.Context, user *models.User, sessionId string) (time.Time, error)
	GetSession(ctx context.Context, sessionId string) (*models.Session, error)
	UpdateSessionExpireTime(ctx context.Context, sessionId string) (time.Time, error)

	AddNewUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	GetRoleIdByName(ctx context.Context, role string) (int, error)
}

type Repositories struct {
	PickupPoint PickupPoint
	Reception   Reception
	Auth        Auth
}

func NewRepositories(db postgres.DBPool) *Repositories {
	return &Repositories{
		PickupPoint: pickupPointRepo.NewPickupPointRepo(db),
		Reception:   receptionRepo.NewReceptionRepo(db),
		Auth:        authRepo.NewAuthRepo(db),
	}
}
