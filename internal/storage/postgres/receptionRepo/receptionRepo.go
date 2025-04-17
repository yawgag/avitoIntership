package receptionRepo

import (
	"context"
	"orderPickupPoint/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReceptionRepo struct {
	pool *pgxpool.Pool
}

func NewReceptionRepo(pool *pgxpool.Pool) *ReceptionRepo {
	return &ReceptionRepo{
		pool: pool,
	}
}

func (r *ReceptionRepo) GetStatusNameById(ctx context.Context, id int) (string, error) {
	query := `select name 
				from receptionStatus
				where id = $1`

	var name string
	err := r.pool.QueryRow(ctx, query, id).Scan(&name)

	return name, err
}

func (r *ReceptionRepo) CreateReception(ctx context.Context, pvzId int) (*models.Reception, error) {
	query := `insert into reception(pvzId)
				select $1
				where not exists(
					select 1
					from reception
					where pvzid = $1 and receptionStatus = 1)
				returning id, receptionStartDateTime, pvzid, statusId`

	outReception := &models.Reception{}
	err := r.pool.QueryRow(ctx, query, pvzId).Scan(&outReception.Id, &outReception.DateTime, &outReception.PickupPointId, &outReception.Status)
	if err != nil {
		return nil, err
	}
	return outReception, nil
}
