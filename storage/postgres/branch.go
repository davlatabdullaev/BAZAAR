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

type branchRepo struct {
	pool *pgxpool.Pool
	log  logger.ILogger
}

func NewBranchRepo(pool *pgxpool.Pool, log logger.ILogger) storage.IBranchRepo {
	return &branchRepo{
		pool: pool,
		log:  log,
	}
}

func (b *branchRepo) Create(ctx context.Context, branch models.CreateBranch) (string, error) {

	id := uuid.New()

	query := `insert into branch (
		id, 
		name, 
		address) 
		values ($1, $2, $3)`

	_, err := b.pool.Exec(ctx, query,
		id,
		branch.Name,
		branch.Address,
	)
	if err != nil {
		b.log.Error("error while inserting branch", logger.Error(err))
		return "", err
	}

	return id.String(), nil
}

func (b *branchRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Branch, error) {

	var updatedAt = sql.NullTime{}

	branch := models.Branch{}

	query := `select
	 id, 
	 name, 
	 address, 
	 created_at, 
	 updated_at
	 from branch where deleted_at is null and id = $1`

	row := b.pool.QueryRow(ctx, query, id.ID)

	err := row.Scan(
		&branch.ID,
		&branch.Name,
		&branch.Address,
		&branch.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		b.log.Error("error while selecting branch", logger.Error(err))
		return models.Branch{}, err
	}

	if updatedAt.Valid {
		branch.UpdatedAt = updatedAt.Time
	}

	return branch, nil

}

func (b *branchRepo) GetList(ctx context.Context, request models.GetListRequest) (models.BranchsResponse, error) {

	var (
		updatedAt         = sql.NullTime{}
		branchs           = []models.Branch{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from branch where deleted_at is null`

	if search != "" {
		countQuery += fmt.Sprintf(` and name ilike '%%%s%%' or address ilike '%%%s%%'`, search, search)
	}
	if err := b.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		b.log.Error("error is while selecting count", logger.Error(err))
		return models.BranchsResponse{}, err
	}

	query = `select 
	id, 
	name, 
	address, 
	created_at, 
	updated_at 
	from branch where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` and name ilike '%%%s%%' or address ilike '%%%s%%'`, search, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := b.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		b.log.Error("error is while selecting branch", logger.Error(err))
		return models.BranchsResponse{}, err
	}

	for rows.Next() {
		branch := models.Branch{}
		if err = rows.Scan(
			&branch.ID,
			&branch.Name,
			&branch.Address,
			&branch.CreatedAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning branch data", logger.Error(err))
			return models.BranchsResponse{}, err
		}

		if updatedAt.Valid {
			branch.UpdatedAt = updatedAt.Time
		}

		branchs = append(branchs, branch)

	}

	return models.BranchsResponse{
		Branchs: branchs,
		Count:   count,
	}, nil
}

func (b *branchRepo) Update(ctx context.Context, request models.UpdateBranch) (string, error) {

	query := `update branch
   set name = $1,
    address = $2, 
	updated_at = $3
   where id = $4  
   `

	_, err := b.pool.Exec(ctx, query,
		request.Name,
		request.Address,
		time.Now(),
		request.ID)
	if err != nil {
		b.log.Error("error while updating branch data...", logger.Error(err))
		return "", err
	}

	return request.ID, nil
}

func (b *branchRepo) Delete(ctx context.Context, id string) error {

	query := `
	update branch
	 set deleted_at = $1
	  where id = $2
	`

	_, err := b.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		b.log.Error("error while deleting branch by id", logger.Error(err))
		return err
	}
	return nil
}
