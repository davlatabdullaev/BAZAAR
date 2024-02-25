package postgres

import (
	"bazaar/api/models"
	"bazaar/pkg/logger"
	"bazaar/storage"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type saleRepo struct {
	pool *pgxpool.Pool
	log  logger.ILogger
}

func NewSaleRepo(pool *pgxpool.Pool, log logger.ILogger) storage.ISaleRepo {
	return &saleRepo{
		pool: pool,
		log:  log,
	}
}

func (s *saleRepo) Create(ctx context.Context, sale models.CreateSale) (string, error) {

	id := uuid.New()

	query := `insert into sale (id, branch_id, shop_assistent_id, cashier_id, payment_type, price, status, client_name) values ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := s.pool.Exec(ctx, query,
		id,
		sale.BranchID,
		sale.ShopAssistantID,
		sale.CashierID,
		sale.PaymentType,
		sale.Price,
		sale.Status,
		sale.ClientName,
	)
	if err != nil {
		s.log.Error("error while inserting sale", logger.Error(err))
		return "", err
	}

	return id.String(), nil
}

func (s *saleRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Sale, error) {

	var updatedAt = sql.NullTime{}

	sale := models.Sale{}

	row := s.pool.QueryRow(ctx, `select id, branch_id, shop_assistent_id, cashier_id, payment_type, price, status, client_name, created_at, updated_at  from sale where deleted_at is null and id = $1`, id.ID)

	err := row.Scan(
		&sale.ID,
		&sale.BranchID,
		&sale.ShopAssistantID,
		&sale.CashierID,
		&sale.PaymentType,
		&sale.Price,
		&sale.Status,
		&sale.ClientName,
		&sale.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		s.log.Error("error while selecting sale", logger.Error(err))
		return models.Sale{}, err
	}

	if updatedAt.Valid {
		sale.UpdatedAt = updatedAt.Time
	}

	return sale, nil
}

func (s *saleRepo) GetList(ctx context.Context, request models.GetListRequest) (models.SalesResponse, error) {

	var (
		updatedAt         = sql.NullTime{}
		sales             = []models.Sale{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from sale where deleted_at is null`

	if search != "" {
		countQuery += fmt.Sprintf(` and status ilike '%%%s%%' or payment_type ilike '%%%s%%'`, search, search)
	}
	if err := s.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", logger.Error(err))
		return models.SalesResponse{}, err
	}

	query = `select 
	id, 
	branch_id, 
	shop_assistent_id, 
	cashier_id, 
	payment_type, 
	price, 
	status, 
	client_name, 
	created_at, 
	updated_at from sale where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` and status ilike '%%%s%%' or payment_type ilike '%%%s%%'`, search, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := s.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting product", logger.Error(err))
		return models.SalesResponse{}, err
	}

	for rows.Next() {
		sale := models.Sale{}
		if err = rows.Scan(
			&sale.ID,
			&sale.BranchID,
			&sale.ShopAssistantID,
			&sale.CashierID,
			&sale.PaymentType,
			&sale.Price,
			&sale.Status,
			&sale.ClientName,
			&sale.CreatedAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning sale data", logger.Error(err))
			return models.SalesResponse{}, err
		}

		if updatedAt.Valid {
			sale.UpdatedAt = updatedAt.Time
		}

		sales = append(sales, sale)

	}

	return models.SalesResponse{
		Sales: sales,
		Count: count,
	}, nil
}

func (s *saleRepo) Update(ctx context.Context, request models.UpdateSale) (string, error) {

	query := `update sale set 
	branch_id = $1, 
	shop_assistent_id = $2,
	cashier_id = $3, 
	payment_type = $4, 
	price = $5, 
	status = $6, 
	client_name = $7, 
	updated_at = $8 
	 where id = $9`

	_, err := s.pool.Exec(ctx, query,
		request.BranchID,
		request.ShopAssistantID,
		request.CashierID,
		request.PaymentType,
		request.Price,
		request.Status,
		request.ClientName,
		time.Now(),
		request.ID,
	)
	if err != nil {
		s.log.Error("error while updating sale data...", logger.Error(err))
		return "", err
	}
	return request.ID, nil
}

func (s *saleRepo) Delete(ctx context.Context, id string) error {

	query := `update sale 
	set deleted_at = $1 
	where id = $2`

	_, err := s.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		s.log.Error("error while deleting sale by id", logger.Error(err))
		return err
	}

	return nil
}

func (s *saleRepo) UpdateSalePrice(ctx context.Context, request models.SaleRequest) (string, error) {

	query := `update sale set 
  price = $1,
  status = $2,
  updated_at = $3 
  where id = $4 `

	if rowsAffected, err := s.pool.Exec(ctx, query, request.TotalPrice, request.Status, time.Now(), request.ID); err != nil {
		if r := rowsAffected.RowsAffected(); r == 0 {
			s.log.Error("error in rows affected ", logger.Error(err))
			return "", err
		}
		s.log.Error("error while updating sale price and status...", logger.Error(err))
		return "", err
	}
	return request.ID, nil

}
