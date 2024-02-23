package postgres

import (
	"bazaar/api/models"
	"bazaar/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncomeProductRepo struct {
	pool *pgxpool.Pool
}

func NewIncomeProductRepo(pool *pgxpool.Pool) storage.IIncomeProductRepo {
	return &IncomeProductRepo{
		pool: pool,
	}
}

func (i *IncomeProductRepo) Create(ctx context.Context, incomeProduct models.CreateIncomeProduct) (string, error) {

	id := uuid.New()

	query := `insert into income_products (
		id,
		income_id,
		product_id,
		price,
		count
	) values ($1, $2, $3, $4, $5)`

	_, err := i.pool.Exec(ctx, query,
		id,
		incomeProduct.IncomeID,
		incomeProduct.ProductID,
		incomeProduct.Price,
		incomeProduct.Count,
	)
	if err != nil {
		log.Println("error while inserting income product", err.Error())
		return "", err
	}

	return id.String(), nil

}

func (i *IncomeProductRepo) Get(ctx context.Context, id models.PrimaryKey) (models.IncomeProduct, error) {

	var updatedAt = sql.NullTime{}

	incomeProduct := models.IncomeProduct{}

	query := `select
	id,
	income_id,
	product_id,
	price,
	count,
	created_at,
	updated_at
	from income_products where deleted_at is null and id = $1
	`

	row := i.pool.QueryRow(ctx, query, id.ID)

	err := row.Scan(
		&incomeProduct.ID,
		&incomeProduct.IncomeID,
		&incomeProduct.ProductID,
		&incomeProduct.Price,
		&incomeProduct.Count,
		&incomeProduct.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		log.Println("error while selecting income products", err.Error())
		return models.IncomeProduct{}, err
	}

	if updatedAt.Valid {
		incomeProduct.UpdatedAt = updatedAt.Time
	}

	return incomeProduct, nil
}

func (i *IncomeProductRepo) GetList(ctx context.Context, request models.GetListRequest) (models.IncomeProductsResponse, error) {

	var (
		updatedAt         = sql.NullTime{}
		incomeProducts    = []models.IncomeProduct{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from income_products where deleted_at is null`

	if search != "" {
		countQuery += fmt.Sprintf(` and product_id ilike '%%%s%%' or price ilike '%%%s%%'`, search, search)
	}
	if err := i.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		log.Println("error while selecting income products count", err.Error())
		return models.IncomeProductsResponse{}, err
	}

	query = `select 
	id,
	income_id,
	product_id,
	price,
	count,
	created_at,
	updated_at
	from income_products where deleted_at is null
	`
	if search != "" {
		countQuery += fmt.Sprintf(` and product_id ilike '%%%s%%' or price ilike '%%%s%%'`, search, search)
	}

	query += `LIMIT $1 OFFSET $2`
	rows, err := i.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		log.Println("error is while selecting income products", err.Error())
		return models.IncomeProductsResponse{}, err
	}

	for rows.Next() {
		incomeProduct := models.IncomeProduct{}
		if err = rows.Scan(
			&incomeProduct.ID,
			&incomeProduct.IncomeID,
			&incomeProduct.ProductID,
			&incomeProduct.Price,
			&incomeProduct.Count,
			&incomeProduct.CreatedAt,
			&updatedAt,
		); err != nil {
			log.Println("error while scanning income_product data", err.Error())
			return models.IncomeProductsResponse{}, err
		}

		if updatedAt.Valid {
			incomeProduct.UpdatedAt = updatedAt.Time
		}

		incomeProducts = append(incomeProducts, incomeProduct)

	}

	return models.IncomeProductsResponse{
		IncomeProducts: incomeProducts,
		Count:          count,
	}, nil
}

func (i *IncomeProductRepo) Update(ctx context.Context, request models.UpdateIncomeProduct) (string, error) {

	query := `update income_products
	 set 
	 income_id = $1,
	 product_id = $2,
	 price = $3,
	 count = $4,
	 updated_at = $5
	 where id = $6
	 `

	_, err := i.pool.Exec(ctx, query,
		request.IncomeID,
		request.ProductID,
		request.Price,
		request.Count,
		time.Now(),
		request.ID,
	)
	if err != nil {
		log.Println("error while updating income product data...", err.Error())
		return "", err
	}

	return request.ID, nil
}

func (i *IncomeProductRepo) Delete(ctx context.Context, id string) error {

	query := `
	update income_products
	 set deleted_at = $1
	  where id = $2
	`

	_, err := i.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting income product by id", err.Error())
		return err
	}

	return nil
}
