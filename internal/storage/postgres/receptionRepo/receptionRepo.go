package receptionRepo

import (
	"context"
	"fmt"
	"orderPickupPoint/internal/models"
	"time"

	"github.com/google/uuid"
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
				from reception_statuses
				where id = $1`

	var name string
	err := r.pool.QueryRow(ctx, query, id).Scan(&name)

	return name, err
}

func (r *ReceptionRepo) GetProductTypeIdByName(ctx context.Context, name string) (int, error) {
	query := `select id
				from product_types
				where name = $1`

	var id int
	err := r.pool.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (r *ReceptionRepo) CreateReception(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error) {
	query := `insert into receptions(pvz_id)
				select $1
				where not exists(
					select 1
					from receptions
					where pvz_id = $1 and status_id = 1)
				returning id, reception_start_datetime, pvz_id, status_id`

	outReception := &models.Reception{}
	err := r.pool.QueryRow(ctx, query, pvzId).Scan(&outReception.Id, &outReception.DateTime, &outReception.PickupPointId, &outReception.StatusId)
	if err != nil {
		return nil, err
	}
	return outReception, nil
}

func (r *ReceptionRepo) AddProductToReception(ctx context.Context, product *models.Product, pvzId uuid.UUID) (*models.Product, error) {
	queryOpenReception := `select id
							from receptions
							where pvz_id = $1 and status_id = 1`

	queryAddProduct := `insert into products(type_id)
						values ($1)
						returning id, added_at`

	query_reception_product := `insert into reception_products(reception_id, product_id)
								values($1, $2)`

	var (
		receptionId uuid.UUID
		productId   uuid.UUID
		addedAt     time.Time
	)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, queryOpenReception, pvzId).Scan(&receptionId)

	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(ctx, queryAddProduct, product.TypeId).Scan(&productId, &addedAt)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(ctx, query_reception_product, receptionId, productId)
	if err != nil {
		return nil, err
	}

	tx.Commit(ctx)

	outReception := &models.Product{
		Id:          productId,
		AddedAt:     addedAt,
		TypeId:      product.TypeId,
		ReceptionId: receptionId,
	}

	return outReception, nil
}

func (r *ReceptionRepo) DeleteLastProductInReception(ctx context.Context, pvzId uuid.UUID) error {
	queryReceptionIsOpen := `select id
							from receptions
							where pvz_id = $1 and status_id = 1`

	queryProductIndex := `select id
							from reception_products rp
							left join products p on p.id = rp.product_id
							where reception_id = $1
							order by p.added_at desc
							limit 1`

	queryDeleteProduct := `delete from products p
							where id = $1`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var (
		receptionId uuid.UUID
		productId   uuid.UUID
	)
	err = tx.QueryRow(ctx, queryReceptionIsOpen, pvzId).Scan(&receptionId)
	if err != nil {
		return err
	}

	err = tx.QueryRow(ctx, queryProductIndex, receptionId).Scan(&productId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, queryDeleteProduct, productId)
	if err != nil {
		return err
	}

	tx.Commit(ctx)

	return nil
}

// TODO: why it's looking so sad? :-(
func (r *ReceptionRepo) CloseReception(ctx context.Context, pvzId uuid.UUID) error {
	query := `update receptions
				set status_id = 2
				where pvz_id = $1 and status_id = 1`
	_, err := r.pool.Exec(ctx, query, pvzId)
	fmt.Println(err)

	return err
}
