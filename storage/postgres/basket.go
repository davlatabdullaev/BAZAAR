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

type basketRepo struct {
	db *sql.DB
}

func NewBasketRepo(db *sql.DB) storage.IBasketRepo {
	return basketRepo{
		db: db,
	}
}

func (b basketRepo) Create(basket models.CreateBasket) (string, error) {

	id := uuid.New()

	query := `insert into category (id, sale_id, product_id, quantity, price, updated_at) values ($1, $2, $3, $4, $5, $6)`

	res, err := b.db.Exec(query,
		id,
		basket.SaleID,
		basket.ProductID,
		basket.Quantity,
		basket.Price,
		check.TimeNow(),
	)
	if err != nil {
		log.Println("error while inserting basket", err.Error())
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

func (b basketRepo) Get(id models.PrimaryKey) (models.Basket, error) {

	basket := models.Basket{}

	row := b.db.QueryRow(`select id, sale_id, product_id, quantity, price, created_at, updated_at from basket where deleted_at is null and id = $1`, id)

	err := row.Scan(
		&basket.ID,
		&basket.SaleID,
		&basket.ProductID,
		&basket.Quantity,
		&basket.Price,
		&basket.CreatedAt,
		&basket.UpdatedAt,
	)

	if err != nil {
		log.Println("error while selecting basket", err.Error())
		return models.Basket{}, err
	}

	return basket, nil
}

func (b basketRepo) GetList(request models.GetListRequest) (models.BasketsResponse, error) {

	var (
		baskets           = []models.Basket{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from basket `

	if search != "" {
		countQuery += fmt.Sprintf(`where price ilike '%%%s%%'`, search)
	}
	if err := b.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
		return models.BasketsResponse{}, err
	}

	query = `select id, sale_id, product_id, quantity, price, created_at, updated_at from basket where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` where price ilike '%%%s%%'`, search)
	}

	query += `LIMIT $1 OFFSET $2`
	rows, err := b.db.Query(query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting basket", err.Error())
		return models.BasketsResponse{}, err
	}

	for rows.Next() {
		basket := models.Basket{}
		if err = rows.Scan(
			&basket.ID,
			&basket.SaleID,
			&basket.ProductID,
			&basket.Quantity,
			&basket.Price,
			&basket.CreatedAt,
			&basket.UpdatedAt,
		); err != nil {
			fmt.Println("error is while scanning basket data", err.Error())
			return models.BasketsResponse{}, err
		}
		baskets = append(baskets, basket)

	}

	return models.BasketsResponse{
		Baskets: baskets,
		Count:   count,
	}, nil
}

func (b basketRepo) Update(request models.UpdateBasket) (string, error) {

	query := `update basket
   set sale_id = $1,
    product_id = $2, 
	quantity = $3,
	price = $4,
	updated_at = $5 
   where id = $6  
   `

	res, err := b.db.Exec(query,
		request.SaleID,
		request.ProductID,
		request.Quantity,
		request.Price,
		check.TimeNow(),
		request.ID)
	if err != nil {
		log.Println("error while updating basket data...", err.Error())
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

func (b basketRepo) Delete(id string) error {

	query := `
	update basket
	 set deleted_at = $1
	  where id = $2
	`

	res, err := b.db.Exec(query, check.TimeNow(), id)
	if err != nil {
		log.Println("error while deleting basket by id", err.Error())
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
