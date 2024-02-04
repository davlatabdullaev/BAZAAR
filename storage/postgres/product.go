package postgres

import (
	"bazaar/api/models"
	"bazaar/pkg/check"
	"bazaar/storage"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepo(pool *pgxpool.Pool) storage.IProductRepo {
	return &productRepo{
		pool: pool,
	}
}

func (p *productRepo) Create(ctx context.Context, product models.CreateProduct) (string, error) {

	id := uuid.New()

	query := `insert into product (id, name, price, barcode, category_id) values ($1, $2, $3, $4, $5)`

	_, err := p.pool.Exec(ctx, query,
		id,
		product.Name,
		product.Price,
		check.GenerateBarCode(),
		product.CategoryID,
	)
	if err != nil {
		log.Println("error while inserting product", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (p *productRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Product, error) {

	product := models.Product{}

	row := p.pool.QueryRow(ctx, `select id, name, price, barcode, category_id, created_at, updated_at  from product where deleted_at is null and id = $1`, id)

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Barcode,
		&product.CategoryID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		log.Println("error while selecting product", err.Error())
		return models.Product{}, err
	}

	return product, nil
}

func (p *productRepo) GetList(ctx context.Context, request models.GetListRequest) (models.ProductsResponse, error) {

	var (
		products          = []models.Product{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from product `

	if search != "" {
		countQuery += fmt.Sprintf(`where name ilike '%%%s%%'`, search)
	}
	if err := p.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
		return models.ProductsResponse{}, err
	}

	query = `select id, name, price, barcode, category_id, created_at, updated_at from product where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` where name ilike '%%%s%%'`, search)
	}

	query += `LIMIT $1 OFFSET $2`
	rows, err := p.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting product", err.Error())
		return models.ProductsResponse{}, err
	}

	for rows.Next() {
		product := models.Product{}
		if err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Barcode,
			&product.CategoryID,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			fmt.Println("error is while scanning product data", err.Error())
			return models.ProductsResponse{}, err
		}

		products = append(products, product)

	}

	return models.ProductsResponse{
		Products: products,
		Count:    count,
	}, nil
}

func (p *productRepo) Update(ctx context.Context, request models.UpdateProduct) (string, error) {

	query := `update product
   set name = $1, price = $2, category_id = $3, updated_at = $4
   where id = $5  
   `

	_, err := p.pool.Exec(ctx, query,
		request.Name,
		request.Price,
		request.CategoryID,
		time.Now(),
		request.ID)
	if err != nil {
		log.Println("error while updating product data...", err.Error())
		return "", err
	}
	return request.ID, nil
}

func (p *productRepo) Delete(ctx context.Context, id string) error {

	query := `
	update product
	 set deleted_at = $1
	  where id = $2
	`

	_, err := p.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting product by id", err.Error())
		return err
	}

	return nil
}
