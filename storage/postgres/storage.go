package postgres

import (
	"bazaar/api/models"
	"bazaar/pkg/check"
	"bazaar/storage"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type storageRepo struct {
	db *sql.DB
}

func NewStorageRepo(db *sql.DB) storage.IStorageRepo {
	return storageRepo{
		db: db,
	}
}

func (s storageRepo) Create(storage models.CreateStorage) (string, error) {

	id := uuid.New()

	query := `insert into storage (id, product_id, branch_id, count, updated_at) values ($1, $2, $3, $4, now())`

	res, err := s.db.Exec(query,
		id,
		storage.ProductID,
		storage.BranchID,
		storage.Count,
	)
	if err != nil {
		log.Println("error while inserting storage", err.Error())
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

	return "", nil
}

func (s storageRepo) Get(id models.PrimaryKey) (models.Storage, error) {

	storage := models.Storage{}

	row := s.db.QueryRow(`select id, product_id, branch_id, count, created_at, updated_at  from storage where deleted_at is null and id = $1`, id)

	err := row.Scan(
		&storage.ID,
		&storage.ProductID,
		&storage.BranchID,
		&storage.Count,
		&storage.CreatedAt,
		&storage.UpdatedAt,
	)

	if err != nil {
		log.Println("error while selecting storage", err.Error())
		return models.Storage{}, err
	}

	return storage, nil
}

func (s storageRepo) GetList(request models.GetListRequest) (models.StoragesResponse, error) {

	var (
		storages          = []models.Storage{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from storage `

	if search != "" {
		countQuery += fmt.Sprintf(`where count ilike '%%%s%%'`, search)
	}
	if err := s.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
		return models.StoragesResponse{}, err
	}

	query = `select id, product_id, branch_id, count, created_at, updated_at from storage where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` where count ilike '%%%s%%'`, search)
	}

	query += `LIMIT $1 OFFSET $2`
	rows, err := s.db.Query(query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting product", err.Error())
		return models.StoragesResponse{}, err
	}

	for rows.Next() {
		storage := models.Storage{}
		if err = rows.Scan(
			&storage.ID,
			&storage.ProductID,
			&storage.BranchID,
			&storage.Count,
			&storage.CreatedAt,
			&storage.UpdatedAt,
		); err != nil {
			fmt.Println("error is while scanning storage data", err.Error())
			return models.StoragesResponse{}, err
		}

		storages = append(storages, storage)

	}

	return models.StoragesResponse{
		Storages: storages,
		Count:    count,
	}, nil
}

func (s storageRepo) Update(request models.UpdateStorage) (string, error) {

	query := `update storage set product_id = $1, branch_id = $2, count =$3, updated_at = $4 where id = $5`

	res, err := s.db.Exec(query,
		request.ProductID,
		request.BranchID,
		request.Count,
		check.TimeNow(),
		request.ID,
	)
	if err != nil {
		log.Println("error while updating storage data...", err.Error())
		return "", err
	}

	rowaAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("error while getting rows affected", err.Error())
		return "", err
	}

	if rowaAffected == 0 {
		log.Println("no rows affected during insert")
		return "", errors.New("no rows affected during insert")
	}

	return "", nil
}

func (s storageRepo) Delete(id string) error {

	query := `update storage 
	set deleted_at = $1 
	where id = $2`

	res, err := s.db.Exec(query, check.TimeNow(), id)
	if err != nil {
		log.Println("error while deleting storage by id", err.Error())
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
