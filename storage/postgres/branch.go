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

type branchRepo struct {
	pool *pgxpool.Pool
}

func NewBranchRepo(pool *pgxpool.Pool) storage.IBranchRepo {
	return &branchRepo{
		pool: pool,
	}
}

func (b *branchRepo) Create(ctx context.Context, branch models.CreateBranch) (string, error) {

	id := uuid.New()

	query := `insert into branch (id, name, address) values ($1, $2, $3)`

	_, err := b.pool.Exec(ctx, query,
		id,
		branch.Name,
		branch.Address,
	)
	if err != nil {
		log.Println("error while inserting branch", err.Error())
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
		log.Println("error while selecting branch", err.Error())
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
		countQuery += fmt.Sprintf(` and name ilike '%%%s%%'`, search)
	}
	if err := b.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
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
		query += fmt.Sprintf(` and name ilike '%%%s%%'`, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := b.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting branch", err.Error())
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
			fmt.Println("error is while scanning branch data", err.Error())
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
		log.Println("error while updating branch data...", err.Error())
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
		log.Println("error while deleting branch by id", err.Error())
		return err
	}
	return nil
}
