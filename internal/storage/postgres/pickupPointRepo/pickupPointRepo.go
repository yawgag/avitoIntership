package pickupPointRepo

import (
	"context"
	"orderPickupPoint/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PickupPointRepo struct {
	pool *pgxpool.Pool
}

func NewPickupPointRepo(pool *pgxpool.Pool) *PickupPointRepo {
	return &PickupPointRepo{
		pool: pool,
	}
}

func (r *PickupPointRepo) GetCityIdByName(ctx context.Context, name string) (int, error) {
	query := `select id
				from cities
				where name = $1`

	var id int
	err := r.pool.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (r *PickupPointRepo) Create(ctx context.Context, pickupPoint *models.PickupPoint) (*models.PickupPoint, error) {
	query := `insert into pvzs(city_id)
				values($1)
				returning id, reg_date`

	outPickupPoint := &models.PickupPoint{}
	err := r.pool.QueryRow(ctx, query, pickupPoint.CityId).Scan(&outPickupPoint.Id, &outPickupPoint.RegDate)
	if err != nil {
		return nil, err
	}

	return outPickupPoint, nil
}
