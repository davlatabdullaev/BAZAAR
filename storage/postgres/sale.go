package postgres

import (
	"bazaar/api/models"
	"bazaar/storage"
	"database/sql"
)

type saleRepo struct {
	db *sql.DB
}

func NewSaleRepo(db *sql.DB) storage.ISaleRepo {
	return saleRepo{
		db: db,
	}
}

func (s saleRepo) Create(product models.CreateSale) (string, error) {

	return "", nil
}

func (s saleRepo) Get(id models.PrimaryKey) (models.Sale, error) {

	return models.Sale{}, nil
}

func (s saleRepo) GetList(request models.GetListRequest) (models.SalesResponse, error) {

	return models.SalesResponse{}, nil
}

func (s saleRepo) Update(request models.UpdateSale) (string, error) {

	return "", nil
}

func (s saleRepo) Delete(id string) error {

	return nil
}
