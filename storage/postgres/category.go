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

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) storage.ICategoryRepo {
	return categoryRepo{
		db: db,
	}
}

func (c categoryRepo) Create(category models.CreateCategory) (string, error) {

	id := uuid.New()

	query := `insert into category (id, name, parent_id) values ($1, $2, $3)`

	res, err := c.db.Exec(query, id, category.Name, category.ParentID)
	if err != nil {
		log.Println("error while inserting category", err.Error())
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

func (c categoryRepo) Get(id models.PrimaryKey) (models.Category, error) {

	category := models.Category{}

	row := c.db.QueryRow(`select id, name, parent_id, created_at, updated_at from category where deleted_at = null and id = $1`, id)

	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.ParentID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		log.Println("error while selecting category", err.Error())
		return models.Category{}, err
	}

	return category, nil
}

func (c categoryRepo) GetList(request models.GetListRequest) (models.CategoriesResponse, error) {
	var (
		categories        = []models.Category{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from category `

	if search != "" {
		countQuery += fmt.Sprintf(`where name ilike '%%%s%%'`, search)
	}
	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
		return models.CategoriesResponse{}, err
	}

	query = `select id, name, parent_id, created_at, updated_at from category where deleted_at = null`

	if search != "" {
		query += fmt.Sprintf(` where name ilike '%%%s%%'`, search)
	}

	query += `LIMIT $1 OFFSET $2`
	rows, err := c.db.Query(query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting category", err.Error())
		return models.CategoriesResponse{}, err
	}

	for rows.Next() {
		category := models.Category{}
		if err = rows.Scan(&category.ID, &category.Name, &category.ParentID, &category.CreatedAt, &category.UpdatedAt); err != nil {
			fmt.Println("error is while scanning category data", err.Error())
			return models.CategoriesResponse{}, err
		}
		categories = append(categories, category)

	}

	return models.CategoriesResponse{
		Categories: categories,
		Count:      count,
	}, nil
}

func (c categoryRepo) Update(request models.UpdateCategory) (string, error) {

	query := `update category
   set name = $1, parent_id = $2, updated_at = $3 
   where id = $4  
   `

	res, err := c.db.Exec(query, request.Name, request.ParentID, check.TimeNow(), request.ID)
	if err != nil {
		log.Println("error while updating category data...", err.Error())
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

func (c categoryRepo) Delete(id string) error {

	query := `
	update category
	 set deleted_at = $1
	  where id = $2
	`

	res, err := c.db.Exec(query, check.TimeNow(), id)
	if err != nil {
		log.Println("error while deleting category by id", err.Error())
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
