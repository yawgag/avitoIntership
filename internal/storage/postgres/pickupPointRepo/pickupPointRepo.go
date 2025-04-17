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
				from city
				where name = $1`

	var id int
	err := r.pool.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (r *PickupPointRepo) Create(ctx context.Context, pickupPoint *models.PickupPoint) (*models.PickupPoint, error) {
	query := `insert into pvz(cityid)
				values($1)
				returning id, regDate`

	cityId, err := r.GetCityIdByName(ctx, pickupPoint.City)
	if err != nil {
		return nil, err
	}

	outPickupPoint := &models.PickupPoint{}
	err = r.pool.QueryRow(ctx, query, cityId).Scan(&outPickupPoint.Id, &outPickupPoint.RegDate)
	if err != nil {
		return nil, err
	}
	outPickupPoint.City = pickupPoint.City

	return outPickupPoint, nil
}
