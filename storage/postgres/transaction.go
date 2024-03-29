package postgres

import (
	"bazaar/api/models"
	"bazaar/pkg/logger"
	"bazaar/storage"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionRepo struct {
	pool *pgxpool.Pool
	log  logger.ILogger
}

func NewTransactionRepo(pool *pgxpool.Pool, log logger.ILogger) storage.ITransactionRepo {
	return &transactionRepo{
		pool: pool,
		log:  log,
	}
}

func (t *transactionRepo) Create(ctx context.Context, request models.CreateTransactions) (string, error) {

	id := uuid.New()

	query := `insert into transactions (
		id, 
		sale_id, 
		staff_id, 
		transaction_type,
		source_type, 
		amount, 
		description) 
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
		t.log.Error("error while inserting transaction data", logger.Error(err))
		return "", err
	}

	return id.String(), nil
}

func (t *transactionRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Transactions, error) {

	var updatedAt = sql.NullTime{}

	transaction := models.Transactions{}

	query := `select 
	id, 
	sale_id, 
	staff_id, 
	transaction_type,
	source_type, 
	amount, 
	description, 
	created_at, 
	updated_at 
	from transactions
	 where deleted_at is null and id = $1`

	row := t.pool.QueryRow(ctx, query, id.ID)

	err := row.Scan(
		&transaction.ID,
		&transaction.SaleID,
		&transaction.StaffID,
		&transaction.TransactionType,
		&transaction.SourceType,
		&transaction.Amount,
		&transaction.Description,
		&transaction.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		t.log.Error("error while selecting transaction data", logger.Error(err))
		return models.Transactions{}, err
	}

	if updatedAt.Valid {
		transaction.UpdatedAt = updatedAt.Time
	}

	return transaction, nil
}

func (t *transactionRepo) GetList(ctx context.Context, request models.GetListTransactionsRequest) (models.TransactionsResponse, error) {
	var (
		updatedAt         = sql.NullTime{}
		page              = request.Page
		offset            = (page - 1) * request.Limit
		transactions      = []models.Transactions{}
		fromAmount        = request.FromAmount
		toAmount          = request.ToAmount
		count             = 0
		query, countQuery string
	)

	countQuery = `select count(1) from transactions where deleted_at is null `
	if fromAmount != 0 && toAmount != 0 {
		countQuery += fmt.Sprintf(` and amount between %f and %f `, fromAmount, toAmount)
	} else if fromAmount != 0 && toAmount == 0 {
		countQuery += ` and amount >= ` + strconv.FormatFloat(fromAmount, 'f', 2, 64)
	} else if toAmount != 0 && fromAmount == 0 {
		countQuery += ` and amount <= ` + strconv.FormatFloat(toAmount, 'f', 2, 64)

	}
	if err := t.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while scanning row", logger.Error(err))
		return models.TransactionsResponse{}, err
	}

	query = `select 
	id, 
	sale_id, 
	staff_id, 
	transaction_type, 
	source_type, 
	amount,
    description, 
	created_at, 
	updated_at from transactions where deleted_at is null `

	if fromAmount != 0 && toAmount != 0 {
		query += fmt.Sprintf(` and amount between %f and %f `, fromAmount, toAmount)
	} else if fromAmount != 0 && toAmount == 0 {
		query += ` and amount >= ` + strconv.FormatFloat(fromAmount, 'f', 2, 64)
	} else if toAmount != 0 && fromAmount == 0 {
		query += ` and amount <= ` + strconv.FormatFloat(toAmount, 'f', 2, 64)

	}

	query += ` LIMIT $1 OFFSET $2 `

	rows, err := t.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting all transaction", logger.Error(err))
		return models.TransactionsResponse{}, err
	}

	for rows.Next() {
		transaction := models.Transactions{}
		if err = rows.Scan(
			&transaction.ID,
			&transaction.SaleID,
			&transaction.StaffID,
			&transaction.TransactionType,
			&transaction.SourceType,
			&transaction.Amount,
			&transaction.Description,
			&transaction.CreatedAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning rows", logger.Error(err))
			return models.TransactionsResponse{}, err
		}

		if updatedAt.Valid {
			transaction.UpdatedAt = updatedAt.Time
		}

		transactions = append(transactions, transaction)
	}
	return models.TransactionsResponse{
		Transactions: transactions,
		Count:        count,
	}, nil
}

func (t *transactionRepo) Update(ctx context.Context, request models.UpdateTransactions) (string, error) {

	query := `update transactions
   set 
   sale_id = $1, 
   staff_id = $2, 
   transaction_type = $3,
   source_type = $4, 
   amount = $5, 
   description = $6, 
   updated_at = $7
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
		t.log.Error("error while updating transaction data...", logger.Error(err))
		return "", err
	}
	return request.ID, nil
}

func (t *transactionRepo) Delete(ctx context.Context, id string) error {

	query := `
	update transactions
	 set deleted_at = $1
	  where id = $2
	`

	_, err := t.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		t.log.Error("error while deleting transaction by id", logger.Error(err))
		return err
	}

	return nil
}

func (t *transactionRepo) UpdateStaffBalanceAndCreateTransaction(ctx context.Context, request models.UpdateStaffBalanceAndCreateTransaction) error {

	transaction, err := t.pool.Begin(ctx)

	defer func() {

		if err != nil {
			transaction.Rollback(ctx)
		} else {
			transaction.Commit(ctx)
		}

	}()

	queryForUpdateStaffBalance := `update staff set
	 balance = balance + $1,
    updated_at = $2
	  where id = $3`

	_, err = transaction.Exec(ctx, queryForUpdateStaffBalance, request.UpdateCashierBalance.Amount, time.Now(), request.UpdateCashierBalance.StaffID)
	if err != nil {
		t.log.Error("error while update staff balance")
		return err
	}

	queryForCreateTransaction := `insert into transactions (
		id, 
		sale_id, 
		staff_id, 
		transaction_type,
		source_type, 
		amount, 
		description) 
	values 
	($1, $2, $3, $4, $5, $6, $7)`

	_, err = transaction.Exec(ctx, queryForCreateTransaction,
		uuid.New().String(),
		request.SaleID,
		request.UpdateCashierBalance.StaffID,
		request.TransactionType,
		request.SourceType,
		request.Amount,
		request.Description,
	)
	if err != nil {
		t.log.Error("error while creating transaction data", logger.Error(err))
		return err
	}

	if request.UpdateShopAssistantBalance.StaffID != "" {

		queryForUpdateStaffBalance := `update staff set
	 balance = balance + $1,
     updated_at = $2
	 where id = $3`

		_, err = transaction.Exec(ctx, queryForUpdateStaffBalance, request.UpdateShopAssistantBalance.Amount, time.Now(), request.UpdateShopAssistantBalance.StaffID)
		if err != nil {
			t.log.Error("error while update staff balance")
			return err
		}

		queryForCreateTransaction := `insert into transactions (
	   id, 
	   sale_id, 
	   staff_id, 
	   transaction_type,
	   source_type, 
	   amount, 
	   description) 
   values 
   ($1, $2, $3, $4, $5, $6, $7)`

		_, err = transaction.Exec(ctx, queryForCreateTransaction,
			uuid.New().String(),
			request.SaleID,
			request.UpdateShopAssistantBalance.StaffID,
			request.TransactionType,
			request.SourceType,
			request.Amount,
			request.Description,
		)
		if err != nil {
			t.log.Error("error while creating transaction data", logger.Error(err))
			return err
		}

	}

	return nil
}
