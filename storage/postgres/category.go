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

type categoryRepo struct {
	pool *pgxpool.Pool
	log  logger.ILogger
}

func NewCategoryRepo(pool *pgxpool.Pool, log logger.ILogger) storage.ICategoryRepo {
	return &categoryRepo{
		pool: pool,
		log:  log,
	}
}

func (c *categoryRepo) Create(ctx context.Context, category models.CreateCategory) (string, error) {

	id := uuid.New()

	query := `insert into category (id, name, parent_id) values ($1, $2, $3)`

	_, err := c.pool.Exec(ctx, query,
		id,
		category.Name,
		category.ParentID,
	)
	if err != nil {
		c.log.Error("error while inserting category", logger.Error(err))
		return "", err
	}

	return id.String(), nil
}

func (c *categoryRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Category, error) {

	var updatedAt = sql.NullTime{}

	category := models.Category{}

	row := c.pool.QueryRow(ctx, `select 
	id, 
	name, 
	parent_id, 
	created_at, 
	updated_at from category where deleted_at is null and id = $1`, id.ID)

	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.ParentID,
		&category.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		c.log.Error("error while selecting category", logger.Error(err))
		return models.Category{}, err
	}

	if updatedAt.Valid {
		category.UpdatedAt = updatedAt.Time
	}

	return category, nil
}

func (c *categoryRepo) GetList(ctx context.Context, request models.GetListRequest) (models.CategoriesResponse, error) {
	var (
		updatedAt         = sql.NullTime{}
		categories        = []models.Category{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from category where deleted_at is null`

	if search != "" {
		countQuery += fmt.Sprintf(` and name ilike '%%%s%%'`, search)
	}
	if err := c.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", logger.Error(err))
		return models.CategoriesResponse{}, err
	}

	query = `select id, name, parent_id, created_at, updated_at from category where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` and name ilike '%%%s%%'`, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := c.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting category", logger.Error(err))
		return models.CategoriesResponse{}, err
	}

	for rows.Next() {
		category := models.Category{}
		var parentID sql.NullString
		if err = rows.Scan(
			&category.ID,
			&category.Name,
			&parentID,
			&category.CreatedAt,
			&updatedAt); err != nil {
			fmt.Println("error is while scanning category data", logger.Error(err))
			return models.CategoriesResponse{}, err
		}

		if parentID.Valid {
			category.ParentID = parentID.String
		} else {
			category.ParentID = ""
		}

		if updatedAt.Valid {
			category.UpdatedAt = updatedAt.Time
		}

		categories = append(categories, category)

	}

	return models.CategoriesResponse{
		Categories: categories,
		Count:      count,
	}, nil
}

func (c *categoryRepo) Update(ctx context.Context, request models.UpdateCategory) (string, error) {

	query := `update category
   set name = $1, parent_id = $2, updated_at = $3 
   where id = $4  
   `

	_, err := c.pool.Exec(ctx, query,
		request.Name,
		request.ParentID,
		time.Now(),
		request.ID)
	if err != nil {
		c.log.Error("error while updating category data...", logger.Error(err))
		return "", err
	}

	return request.ID, nil
}

func (c *categoryRepo) Delete(ctx context.Context, id string) error {

	query := `
	update category
	 set deleted_at = $1
	  where id = $2
	`

	_, err := c.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		c.log.Error("error while deleting category by id", logger.Error(err))
		return err
	}
	return nil
}
