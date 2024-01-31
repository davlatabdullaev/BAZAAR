package postgres

import (
	"bazaar/api/models"
	"bazaar/storage"
	"database/sql"
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

	return "", nil
}

func (b basketRepo) Get(id models.PrimaryKey) (models.Basket, error) {

	return models.Basket{}, nil
}

func (b basketRepo) GetList(request models.GetListRequest) (models.BasketsResponse, error) {

	return models.BasketsResponse{}, nil
}

func (b basketRepo) Update(request models.UpdateBasket) (string, error) {

	return "", nil
}

func (b basketRepo) Delete(id string) error {

	return nil
}
