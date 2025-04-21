package pickupPointRepo

import (
	"context"
	"fmt"
	"orderPickupPoint/internal/models"
	"orderPickupPoint/internal/storage/postgres"
	"strconv"
)

type PickupPointRepo struct {
	pool postgres.DBPool
}

func NewPickupPointRepo(pool postgres.DBPool) *PickupPointRepo {
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

func (r *PickupPointRepo) GetFilteredInfo(ctx context.Context, filter *models.PvzFilter) ([]models.PvzFilteredInfo, error) {
	queryData := []interface{}{}

	query := `select 	p.id, 
						c.name, 
						p.reg_date, 
						r.id,
						r.reception_start_datetime, 
						prod.id, 
						prod.added_at, 
						pt.name 
				from pvzs p 
				join receptions r on p.id = r.pvz_id
				join reception_products rp on rp.reception_id = r.id
				join products prod on prod.id = rp.product_id
				join product_types pt on pt.id = prod.type_id
				join cities c on p.city_id = c.id`

	if filter.EndDate != nil && filter.StartDate != nil {
		query += "\nwhere prod.added_at between $1 and $2"
		fmt.Println(filter.StartDate, filter.EndDate)
		queryData = append(queryData, filter.StartDate, filter.EndDate)
	}

	query += fmt.Sprintf("\norder by r.reception_start_datetime\nlimit $%s offset $%s;", strconv.Itoa(len(queryData)+1), strconv.Itoa(len(queryData)+2))

	offset := filter.PageLimit * (filter.Page - 1)

	queryData = append(queryData, filter.PageLimit, offset)
	fmt.Println(queryData...)
	rows, err := r.pool.Query(ctx, query, queryData...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.PvzFilteredInfo
	for rows.Next() {
		var row models.PvzFilteredInfo
		err := rows.Scan(
			&row.PvzID,
			&row.CityName,
			&row.RegDate,
			&row.ReceptionID,
			&row.ReceptionTime,
			&row.ProductID,
			&row.AddedAt,
			&row.ProductType,
		)
		if err != nil {
			continue
		}
		out = append(out, row)
	}
	return out, nil

}
