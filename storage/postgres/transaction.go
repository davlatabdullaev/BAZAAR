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

type transactionRepo struct {
	pool *pgxpool.Pool
}

func NewTransactionRepo(pool *pgxpool.Pool) storage.ITransactionRepo {
	return &transactionRepo{
		pool: pool,
	}
}

func (t *transactionRepo) Create(ctx context.Context, request models.CreateTransaction) (string, error) {

	id := uuid.New()

	query := `insert into transaction (id, sale_id, staff_id, transaction_type,
		source_type, amount, description) 
	values 
	($1, $2, $3, $4, $5, $6, $7)`

	_, err := t.pool.Exec(ctx, query,
		id,
		request.SaleID,
		request.StaffID,
		request.TransactionType,
		request.SourceType,
		request.Amount,
		request.Description,
	)
	if err != nil {
		log.Println("error while inserting transaction data", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (t *transactionRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Transaction, error) {

	transaction := models.Transaction{}

	query := `select id, sale_id, staff_id, transaction_type,
	source_type, amount, description, 
	 created_at, updated_at from tarif
	 where deleted_at is null and id = $1`

	row := t.pool.QueryRow(ctx, query, id)

	err := row.Scan(
		&transaction.ID,
		&transaction.SaleID,
		&transaction.StaffID,
		&transaction.TransactionType,
		&transaction.SourceType,
		&transaction.Amount,
		&transaction.Description,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		log.Println("error while selecting transaction data", err.Error())
		return models.Transaction{}, err
	}

	return transaction, nil
}

func (t *transactionRepo) GetList(ctx context.Context, request models.GetListRequest) (models.TransactionsResponse, error) {

	var (
		transactions      = []models.Transaction{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from transaction `

	if search != "" {
		countQuery += fmt.Sprintf(`where description ilike '%%%s%%'`, search)
	}
	if err := t.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting transaction count", err.Error())
		return models.TransactionsResponse{}, err
	}

	query = `select id, sale_id, staff_id, 
	transaction_type, source_type,
	amount, description,
	created_at, updated_at
	from tarif 
	where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` where description ilike '%%%s%%'`, search)
	}

	query += `LIMIT $1 OFFSET $2`
	rows, err := t.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting transaction", err.Error())
		return models.TransactionsResponse{}, err
	}

	for rows.Next() {
		transaction := models.Transaction{}
		if err = rows.Scan(
			&transaction.ID,
			&transaction.SaleID,
			&transaction.StaffID,
			&transaction.TransactionType,
			&transaction.SourceType,
			&transaction.Amount,
			&transaction.Description,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		); err != nil {
			fmt.Println("error is while scanning transaction data", err.Error())
			return models.TransactionsResponse{}, err
		}

		transactions = append(transactions, transaction)

	}

	return models.TransactionsResponse{
		Transactions: transactions,
		Count:        count,
	}, nil
}

func (t *transactionRepo) Update(ctx context.Context, request models.UpdateTransaction) (string, error) {

	query := `update transaction
   set sale_id = $1, staff_id = $2, transaction_type = $3,
   source_type = $4, amount = $5, description = $6, updated_at = $7
   where id = $8
   `
	_, err := t.pool.Exec(ctx, query,
		request.SaleID,
		request.StaffID,
		request.TransactionType,
		request.SourceType,
		request.Amount,
		request.Description,
		time.Now(),
		request.ID,
	)
	if err != nil {
		log.Println("error while updating transaction data...", err.Error())
		return "", err
	}
	return request.ID, nil
}

func (t *transactionRepo) Delete(ctx context.Context, id string) error {

	query := `
	update transaction
	 set deleted_at = $1
	  where id = $2
	`

	_, err := t.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting transaction by id", err.Error())
		return err
	}

	return nil
}
