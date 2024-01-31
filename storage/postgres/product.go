package postgres

import (
	"bazaar/api/models"
	"bazaar/storage"
	"database/sql"
)

type productRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) storage.IProductRepo {
	return productRepo{
		db: db,
	}
}

func (p productRepo) Create(product models.CreateProduct) (string, error) {

	return "", nil
}

func (p productRepo) Get(id models.PrimaryKey) (models.Product, error) {

	return models.Product{}, nil
}

func (p productRepo) GetList(request models.GetListRequest) (models.ProductsResponse, error) {

	return models.ProductsResponse{}, nil
}

func (p productRepo) Update(request models.UpdateProduct) (string, error) {

	return "", nil
}

func (p productRepo) Delete(id string) error {

	return nil
}
