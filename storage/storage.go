package storage

import "bazaar/api/models"

type IStorage interface {
	CloseDB()
	Category() ICategoryRepo
	Staff() IStaffRepo
	StorageTransaction() IStorageTransactionRepo
	Tarif() ITarifRepo
	Transaction() ITransactionRepo
	Basket() IBasketRepo
	Branch() IBranchRepo
	Product() IProductRepo
	Sale() ISaleRepo
	Storage() IStorageRepo
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

type IBasketRepo interface {
	Create(models.CreateBasket) (string, error)
	Get(models.PrimaryKey) (models.Basket, error)
	GetList(models.GetListRequest) (models.BasketsResponse, error)
	Update(models.UpdateBasket) (string, error)
	Delete(string) error
}

type IBranchRepo interface {
	Create(models.CreateBranch) (string, error)
	Get(models.PrimaryKey) (models.Branch, error)
	GetList(models.GetListRequest) (models.BranchsResponse, error)
	Update(models.UpdateBranch) (string, error)
	Delete(string) error
}

type IProductRepo interface {
	Create(models.CreateProduct) (string, error)
	Get(models.PrimaryKey) (models.Product, error)
	GetList(models.GetListRequest) (models.ProductsResponse, error)
	Update(models.UpdateProduct) (string, error)
	Delete(string) error
}

type ISaleRepo interface {
	Create(models.CreateSale) (string, error)
	Get(models.PrimaryKey) (models.Sale, error)
	GetList(models.GetListRequest) (models.SalesResponse, error)
	Update(models.UpdateSale) (string, error)
	Delete(string) error
}

type IStorageRepo interface {
	Create(models.CreateStorage) (string, error)
	Get(models.PrimaryKey) (models.Storage, error)
	GetList(models.GetListRequest) (models.StoragesResponse, error)
	Update(models.UpdateStorage) (string, error)
	Delete(string) error
}