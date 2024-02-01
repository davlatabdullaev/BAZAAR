package postgres

import (
	"bazaar/api/models"
	"bazaar/pkg/check"
	"bazaar/storage"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type storageTransactionRepo struct {
	db *sql.DB
}

func NewStorageTransactionRepo(db *sql.DB) storage.IStorageTransactionRepo {
	return storageTransactionRepo{
		db: db,
	}
}

func (s storageTransactionRepo) Create(request models.CreateStorageTransaction) (string, error) {

	id := uuid.New()

	updatedAt := time.Now()

	query := `insert into storage_transaction (id, staff_id, product_id, storage_tranaction_type, price, quantity, updated_at) values ($1, $2, $3, $4, $5, $6, $7)`

	res, err := s.db.Exec(query,
		id,
		request.StaffID,
		request.ProductID,
		request.StorageTransactionType,
		request.Price,
		request.Quantity,
	    updatedAt,
	)
	if err != nil {
		log.Println("error while inserting storage transaction data", err.Error())
		return "", err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("error getting rows affected", err.Error())
		return "", err
	}

	if rowsAffected == 0 {
		log.Println("no rows affected during insert")
		return "", errors.New("no rows affected during insert")
	}

	return id.String(), nil
}

func (s storageTransactionRepo) Get(id models.PrimaryKey) (models.StorageTransaction, error) {

	storageTransaction := models.StorageTransaction{}

	query := `select id, staff_id, product_id, storage_transaction_type, 
	price, quantity, created_at, updated_at from storage_transaction
	 where deleted_at is null and id = $1`

	row := s.db.QueryRow(query, id)

	err := row.Scan(
		&storageTransaction.ID,
		&storageTransaction.StaffID,
		&storageTransaction.ProductID,
		&storageTransaction.StorageTransactionType,
		&storageTransaction.Price,
		&storageTransaction.Quantity,
		&storageTransaction.CreatedAt,
		&storageTransaction.UpdatedAt,
	)

	if err != nil {
		log.Println("error while selecting storage transaction data", err.Error())
		return models.StorageTransaction{}, err
	}

	return storageTransaction, nil
}

func (s storageTransactionRepo) GetList(request models.GetListRequest) (models.StorageTransactionsResponse, error) {

	var (
		storageTransactions = []models.StorageTransaction{}
		count               = 0
		query, countQuery   string
		page                = request.Page
		offset              = (page - 1) * request.Limit
		search              = request.Search
	)

	countQuery = `select count(1) from storage_transaction `

	if search != "" {
		countQuery += fmt.Sprintf(`where storage_transaction_type ilike '%%%s%%'`, search)
	}
	if err := s.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting storage_transaction count", err.Error())
		return models.StorageTransactionsResponse{}, err
	}

	query = `select id, staff_id, product_id, storafe_transaction_type, 
	price, quantity, created_at, updated_at
	from storage_transaction 
	where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` where storage_transaction_type ilike '%%%s%%'`, search)
	}

	query += `LIMIT $1 OFFSET $2`
	rows, err := s.db.Query(query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting storage transaction", err.Error())
		return models.StorageTransactionsResponse{}, err
	}

	for rows.Next() {
		storageTransaction := models.StorageTransaction{}
		if err = rows.Scan(
			&storageTransaction.ID,
			&storageTransaction.StaffID,
			&storageTransaction.ProductID,
			&storageTransaction.StorageTransactionType,
			&storageTransaction.Price,
			&storageTransaction.Quantity,
			&storageTransaction.CreatedAt,
			&storageTransaction.UpdatedAt,
		); err != nil {
			fmt.Println("error is while scanning storage transaction data", err.Error())
			return models.StorageTransactionsResponse{}, err
		}

		storageTransactions = append(storageTransactions, storageTransaction)

	}

	return models.StorageTransactionsResponse{
		StorageTransactions: storageTransactions,
		Count:               count,
	}, nil
}

func (s storageTransactionRepo) Update(request models.UpdateStorageTransaction) (string, error) {

	query := `update storage_transaction
   set staff_id = $1, product_id = $2, storage_transaction_type = $3,
   price = $4, quantity = $5, updated_at = $6
   where id = $7
   `
	res, err := s.db.Exec(query,
	request.StaffID,
	request.ProductID,
	request.StorageTransactionType,
	request.Price,
	request.Quantity,	
	check.TimeNow(),
	)
	if err != nil {
		log.Println("error while updating storage_transaction data...", err.Error())
		return "", err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("error getting rows affected", err.Error())
		return "", err
	}

	if rowsAffected == 0 {
		log.Println("no rows affected during insert")
		return "", errors.New("no rows affected during insert")
	}

	return request.ID, nil
}

func (s storageTransactionRepo) Delete(id string) error {

	query := `
	update storage_transaction
	 set deleted_at = $1
	  where id = $2
	`

	res, err := s.db.Exec(query, check.TimeNow(), id)
	if err != nil {
		log.Println("error while deleting storage_transaction by id", err.Error())
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("error getting rows affected", err.Error())
		return err
	}

	if rowsAffected == 0 {
		log.Println("no rows affected during insert")
		return errors.New("no rows affected during insert")
	}


	return nil
}
