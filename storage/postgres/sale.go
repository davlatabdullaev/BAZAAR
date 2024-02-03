package postgres

import (
	"bazaar/api/models"
	"bazaar/storage"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type saleRepo struct {
	pool *pgxpool.Pool
}

func NewSaleRepo(pool *pgxpool.Pool) storage.ISaleRepo {
	return &saleRepo{
		pool: pool,
	}
}

func (s *saleRepo) Create(ctx context.Context, sale models.CreateSale) (string, error) {

	id := uuid.New()

	query := `insert into sale (id, branch_id, shop_assistent_id, chashier_id, payment_type, price, status, client_name) values ($1, $2, $3, $4, $5, $6, $7, $8)`

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
		log.Println("error while inserting sale", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (s *saleRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Sale, error) {

	sale := models.Sale{}

	row := s.pool.QueryRow(ctx, `select id, branch_id, shop_assistent_id, cashier_id, payment_type, price, status, client_name, created_at, updated_at  from sale where deleted_at is null and id = $1`, id)

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
		&sale.UpdatedAt,
	)

	if err != nil {
		log.Println("error while selecting sale", err.Error())
		return models.Sale{}, err
	}

	return sale, nil
}

func (s *saleRepo) GetList(ctx context.Context, request models.GetListRequest) (models.SalesResponse, error) {

	var (
		sales             = []models.Sale{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from sale `

	if search != "" {
		countQuery += fmt.Sprintf(`where price ilike '%%%s%%'`, search)
	}
	if err := s.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
		return models.SalesResponse{}, err
	}

	query = `select id, branch_id, shop_assistent_id, cashier_id, payment_type, price, status, client_name, created_at, updated_at from sale where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` where price ilike '%%%s%%'`, search)
	}

	query += `LIMIT $1 OFFSET $2`
	rows, err := s.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting product", err.Error())
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
			&sale.UpdatedAt,
		); err != nil {
			fmt.Println("error is while scanning sale data", err.Error())
			return models.SalesResponse{}, err
		}

		sales = append(sales, sale)

	}

	return models.SalesResponse{
		Sales: sales,
		Count: count,
	}, nil
}

func (s *saleRepo) Update(ctx context.Context, request models.UpdateSale) (string, error) {

	query := `update sale set branch_id = $1, shop_assistent_id = $2,
	 cashier_id = $3, payment_type = $4, price = $5, 
	 status = $6, client_name = $7, updated_at = $8 
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
		log.Println("error while updating sale data...", err.Error())
		return "", err
	}
	return "", nil
}

func (s *saleRepo) Delete(ctx context.Context, id string) error {

	query := `update sale 
	set deleted_at = $1 
	where id = $2`

	_, err := s.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting sale by id", err.Error())
		return err
	}

	return nil
}
