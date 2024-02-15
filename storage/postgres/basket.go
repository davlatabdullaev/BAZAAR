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

type basketRepo struct {
	pool *pgxpool.Pool
}

func NewBasketRepo(pool *pgxpool.Pool) storage.IBasketRepo {
	return &basketRepo{
		pool: pool,
	}
}

func (b *basketRepo) Create(ctx context.Context, basket models.CreateBasket) (string, error) {

	id := uuid.New()

	query := `insert into basket (id, sale_id, product_id, quantity, price) values ($1, $2, $3, $4, $5)`

	_, err := b.pool.Exec(ctx, query,
		id,
		basket.SaleID,
		basket.ProductID,
		basket.Quantity,
		basket.Price,
	)
	if err != nil {
		log.Println("error while inserting basket", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (b *basketRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Basket, error) {

	var updatedAt = sql.NullTime{}

	basket := models.Basket{}

	row := b.pool.QueryRow(ctx, `select 
	id,
    sale_id,
    product_id, 
	quantity, 
	price,
    created_at, 
	updated_at
	from basket where deleted_at is null and id = $1`, id.ID)

	err := row.Scan(
		&basket.ID,
		&basket.SaleID,
		&basket.ProductID,
		&basket.Quantity,
		&basket.Price,
		&basket.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		log.Println("error while selecting basket", err.Error())
		return models.Basket{}, err
	}

	if updatedAt.Valid {
		basket.UpdatedAt = updatedAt.Time
	}

	return basket, nil
}

func (b *basketRepo) GetList(ctx context.Context, request models.GetBasketsListRequest) (models.BasketsResponse, error) {

	var (
		updatedAt         = sql.NullTime{}
		baskets           = []models.Basket{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from basket where deleted_at is null `

	if search != "" {
		countQuery += fmt.Sprintf(` and product_id = '%s' or sale_id = '%s'`, search, search)
	}
	if err := b.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
		return models.BasketsResponse{}, err
	}

	query = `select 
	id, 
	sale_id, 
	product_id, 
	quantity, 
	price, 
	created_at, 
	updated_at
	from basket where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` and product_id = '%s' or sale_id = '%s'`, search, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := b.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting basket", err.Error())
		return models.BasketsResponse{}, err
	}

	for rows.Next() {
		basket := models.Basket{}
		if err = rows.Scan(
			&basket.ID,
			&basket.SaleID,
			&basket.ProductID,
			&basket.Quantity,
			&basket.Price,
			&basket.CreatedAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning basket data", err.Error())
			return models.BasketsResponse{}, err
		}

		if updatedAt.Valid {
			basket.UpdatedAt = updatedAt.Time
		}

		baskets = append(baskets, basket)

	}

	return models.BasketsResponse{
		Baskets: baskets,
		Count:   count,
	}, nil
}

func (b *basketRepo) Update(ctx context.Context, request models.UpdateBasket) (string, error) {

	query := `update basket
   set sale_id = $1,
    product_id = $2, 
	quantity = $3,
	price = $4,
	updated_at = $5 
   where id = $6 
   `

	_, err := b.pool.Exec(ctx, query,
		request.SaleID,
		request.ProductID,
		request.Quantity,
		request.Price,
		time.Now(),
		request.ID)
	if err != nil {
		log.Println("error while updating basket data...", err.Error())
		return "", err
	}

	return request.ID, nil
}

func (b *basketRepo) Delete(ctx context.Context, id string) error {

	query := `
	update basket
	 set deleted_at = $1
	  where id = $2
	`

	_, err := b.pool.Exec(ctx,
		query,
		time.Now(),
		id)
	if err != nil {
		log.Println("error while deleting basket by id", err.Error())
		return err
	}
	return nil
}


func (b *basketRepo) UpdateBasketQuantity(ctx context.Context, request models.UpdateBasketQuantity) (string, error) {

	query := `update basket
   set 
    quantity = quantity + $1,
    updated_at = $2 
   where id = $3 
   `

	_, err := b.pool.Exec(ctx, query,
		request.Quantity,
		time.Now(),
		request.ID)
	if err != nil {
		log.Println("error while updating basket quantity...", err.Error())
		return "", err
	}

	return request.ID, nil
}
