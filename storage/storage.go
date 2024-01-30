package storage

import "bazaar/api/models"

type IStorage interface {
	CloseDB()
	Category() ICategoryRepo
	Staff() IStaffRepo
	StorageTransaction() IStorageTransactionRepo
	Tarif() ITarifRepo
	Transaction() ITransactionRepo
}

type ICategoryRepo interface {
	Create(models.CreateCategory) (string, error)
	Get(models.PrimaryKey) (models.Category, error)
	GetList(models.GetListRequest) (models.CategoriesResponse, error)
	Update(models.UpdateCategory) (string, error)
	Delete(string) error
}

type IStaffRepo interface {
	Create(models.CreateStaff) (string, error)
	Get(models.PrimaryKey) (models.Staff, error)
	GetList(models.GetListRequest) (models.StaffsResponse, error)
	Update(models.UpdateStaff) (string, error)
	Delete(string) error
}

type IStorageTransactionRepo interface {
	Create(models.CreateStorageTransaction) (string, error)
	Get(models.PrimaryKey) (models.StorageTransaction, error)
	GetList(models.GetListRequest) (models.StorageTransactionsResponse, error)
	Update(models.UpdateStorageTransaction) (string, error)
	Delete(string) error
}

type ITarifRepo interface {
	Create(models.CreateTarif) (string, error)
	Get(models.PrimaryKey) (models.Tarif, error)
	GetList(models.GetListRequest) (models.TarifsResponse, error)
	Update(models.UpdateTarif) (string, error)
	Delete(string) error
}

type ITransactionRepo interface {
	Create(models.CreateTransaction) (string, error)
	Get(models.PrimaryKey) (models.Transaction, error)
	GetList(models.GetListRequest) (models.TransactionsResponse, error)
	Update(models.UpdateTransaction) (string, error)
	Delete(string) error
}
