package postgres

import (
	"bazaar/api/models"
	"bazaar/storage"
	"database/sql"
)

type storageRepo struct {
	db *sql.DB
}

func NewStorageRepo(db *sql.DB) storage.IStorageRepo {
	return storageRepo{
		db: db,
	}
}

func (s storageRepo) Create(product models.CreateStorage) (string, error) {

	return "", nil
}

func (s storageRepo) Get(id models.PrimaryKey) (models.Storage, error) {

	return models.Storage{}, nil
}

func (s storageRepo) GetList(request models.GetListRequest) (models.StoragesResponse, error) {

	return models.StoragesResponse{}, nil
}

func (s storageRepo) Update(request models.UpdateStorage) (string, error) {

	return "", nil
}

func (s storageRepo) Delete(id string) error {

	return nil
}
