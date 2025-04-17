package storage

import (
	"context"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/storage/postgres/authRepo"
	"orderPickupPoint/internal/storage/postgres/pickupPointRepo"
	"orderPickupPoint/internal/storage/postgres/receptionRepo"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PickupPoint interface {
	// Create(ctx context.Context, data *models.PickupPoint) error
	// ??? GetAllInfo
	Create(ctx context.Context, pickupPoint *models.PickupPoint) (*models.PickupPoint, error)
}

type Reception interface {
	// Create(ctx context.Context, ... ) error
	// AddProduct(ctx context.Context, ...) error
	// DeleteProduct(ctx context.Context, ...) error
	// Close(ctx context.Context, ...) error
	GetStatusNameById(ctx context.Context, id int) (string, error)
	CreateReception(ctx context.Context, pvzId int) (*models.Reception, error)
}

type Auth interface {
	CreateSession(ctx context.Context, user *models.User, sessionId string) (time.Time, error)
	GetSession(ctx context.Context, sessionId string) (*models.Session, error)
	UpdateSessionExpireTime(ctx context.Context, sessionId string) (time.Time, error)

	AddNewUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	GetRoleIdByName(ctx context.Context, role string) (int, error)

	// LOGIN - GetSession(ctx context.Context, sessionId string) (*models.Session, error)
	// LOGOUT - DeleteSession(ctx context.Context, sessionId string) error
}

type Repositories struct {
	PickupPoint PickupPoint
	Reception   Reception
	Auth        Auth
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		PickupPoint: pickupPointRepo.NewPickupPointRepo(db),
		Reception:   receptionRepo.NewReceptionRepo(db),
		Auth:        authRepo.NewAuthRepo(db),
	}
}
