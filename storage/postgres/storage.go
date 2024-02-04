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

type storageRepo struct {
	pool *pgxpool.Pool
}

func NewStorageRepo(pool *pgxpool.Pool) storage.IStorageRepo {
	return &storageRepo{
		pool: pool,
	}
}

func (s *storageRepo) Create(ctx context.Context, storage models.CreateStorage) (string, error) {

	id := uuid.New()

	query := `insert into storage (id, product_id, branch_id, count) values ($1, $2, $3, $4)`

	_, err := s.pool.Exec(ctx, query,
		id,
		storage.ProductID,
		storage.BranchID,
		storage.Count,
	)
	if err != nil {
		log.Println("error while inserting storage", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (s *storageRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Storage, error) {

	storage := models.Storage{}

	row := s.pool.QueryRow(ctx, `select id, product_id, branch_id, count, created_at, updated_at  from storage where deleted_at is null and id = $1`, id.ID)

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

func (s *storageRepo) GetList(ctx context.Context, request models.GetListRequest) (models.StoragesResponse, error) {

	var (
		storages          = []models.Storage{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from storage where deleted_at is null `

	if search != "" {
		countQuery += fmt.Sprintf(`and count ilike '%%%s%%'`, search)
	}
	if err := s.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
		return models.StoragesResponse{}, err
	}

	query = `select id, product_id, branch_id, count, created_at, updated_at from storage where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` where count ilike '%%%s%%'`, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := s.pool.Query(ctx, query, request.Limit, offset)
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

func (s *storageRepo) Update(ctx context.Context, request models.UpdateStorage) (string, error) {

	query := `update storage set product_id = $1, branch_id = $2, count =$3, updated_at = $4 where id = $5`

	_, err := s.pool.Exec(ctx, query,
		request.ProductID,
		request.BranchID,
		request.Count,
		time.Now(),
		request.ID,
	)
	if err != nil {
		log.Println("error while updating storage data...", err.Error())
		return "", err
	}

	return request.ID, nil
}

func (s *storageRepo) Delete(ctx context.Context, id string) error {

	query := `update storage 
	set deleted_at = $1 
	where id = $2`

	_, err := s.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting storage by id", err.Error())
		return err
	}

	return nil
}
