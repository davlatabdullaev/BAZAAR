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

type IncomeRepo struct {
	pool *pgxpool.Pool
	log  logger.ILogger
}

func NewIncomeRepo(pool *pgxpool.Pool, log logger.ILogger) storage.IIncomeRepo {
	return &IncomeRepo{
		pool: pool,
		log:  log,
	}
}

func (i *IncomeRepo) Create(ctx context.Context, income models.CreateIncome) (string, error) {

	id := uuid.New()

	query := `insert into income (
		id,
		branch_id,
		price
	) values ($1, $2, $3)`

	_, err := i.pool.Exec(ctx, query,
		id,
		income.BranchID,
		income.Price,
	)
	if err != nil {
		i.log.Error("error while inserting income", logger.Error(err))
		return "", err
	}

	return id.String(), nil
}

func (i *IncomeRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Income, error) {

	var updatedAt = sql.NullTime{}

	income := models.Income{}

	query := `select
	id,
	branch_id,
	price,
	created_at,
	updated_at
	from income where deleted_at is null and id = $1
	`

	row := i.pool.QueryRow(ctx, query, id.ID)

	err := row.Scan(
		&income.ID,
		&income.BranchID,
		&income.Price,
		&income.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		i.log.Error("error while selecting income", logger.Error(err))
		return models.Income{}, err
	}

	if updatedAt.Valid {
		income.UpdatedAt = updatedAt.Time
	}

	return income, nil
}

func (i *IncomeRepo) GetList(ctx context.Context, request models.GetListRequest) (models.IncomesResponse, error) {

	var (
		updatedAt         = sql.NullTime{}
		incomes           = []models.Income{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from income where deleted_at is null`

	if search != "" {
		countQuery += fmt.Sprintf(` and branch_id ilike '%%%s%%' or price ilike '%%%s%%'`, search, search)
	}
	if err := i.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		i.log.Error("error while selecting count", logger.Error(err))
		return models.IncomesResponse{}, err
	}

	query = `select 
	id,
	branch_id,
	price,
	created_at,
	updated_at
	from income where deleted_at is null
	`
	if search != "" {
		countQuery += fmt.Sprintf(` and branch_id ilike '%%%s%%' or price ilike '%%%s%%'`, search, search)
	}

	query += `LIMIT $1 OFFSET $2`
	rows, err := i.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		i.log.Error("error is while selecting income", logger.Error(err))
		return models.IncomesResponse{}, err
	}

	for rows.Next() {
		income := models.Income{}
		if err = rows.Scan(
			&income.ID,
			&income.BranchID,
			&income.Price,
			&income.CreatedAt,
			&updatedAt,
		); err != nil {
			i.log.Error("error while scanning income data", logger.Error(err))
			return models.IncomesResponse{}, err
		}

		if updatedAt.Valid {
			income.UpdatedAt = updatedAt.Time
		}

		incomes = append(incomes, income)

	}

	return models.IncomesResponse{
		Incomes: incomes,
		Count:   count,
	}, nil
}

func (i *IncomeRepo) Update(ctx context.Context, request models.UpdateIncome) (string, error) {

	query := `update income 
	 set 
	 branch_id = $1,
	 price = $2,
	 updated_at = $3
	 where id = $4
	 `

	_, err := i.pool.Exec(ctx, query,
		request.BranchID,
		request.Price,
		time.Now(),
		request.ID,
	)
	if err != nil {
		i.log.Error("error while updating income data...", logger.Error(err))
		return "", err
	}

	return request.ID, nil
}

func (i *IncomeRepo) Delete(ctx context.Context, id string) error {

	query := `
	update income
	 set deleted_at = $1
	  where id = $2
	`

	_, err := i.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		i.log.Error("error while deleting income by id", logger.Error(err))
		return err
	}

	return nil
}
