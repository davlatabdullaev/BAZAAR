package storage

import (
	"bazaar/api/models"
	"context"
)

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
	Create(context.Context, models.CreateCategory) (string, error)
	Get(context.Context, models.PrimaryKey) (models.Category, error)
	GetList(context.Context, models.GetListRequest) (models.CategoriesResponse, error)
	Update(context.Context, models.UpdateCategory) (string, error)
	Delete(context.Context, string) error
}

type IStaffRepo interface {
	Create(context.Context, models.CreateStaff) (string, error)
	Get(context.Context, models.PrimaryKey) (models.Staff, error)
	GetList(context.Context, models.GetListRequest) (models.StaffsResponse, error)
	Update(context.Context, models.UpdateStaff) (string, error)
	Delete(context.Context, string) error
}

type IStorageTransactionRepo interface {
	Create(context.Context, models.CreateStorageTransaction) (string, error)
	Get(context.Context, models.PrimaryKey) (models.StorageTransaction, error)
	GetList(context.Context, models.GetListRequest) (models.StorageTransactionsResponse, error)
	Update(context.Context, models.UpdateStorageTransaction) (string, error)
	Delete(context.Context, string) error
}

type ITarifRepo interface {
	Create(context.Context, models.CreateTarif) (string, error)
	Get(context.Context, models.PrimaryKey) (models.Tarif, error)
	GetList(context.Context, models.GetListRequest) (models.TarifsResponse, error)
	Update(context.Context, models.UpdateTarif) (string, error)
	Delete(context.Context, string) error
}

type ITransactionRepo interface {
	Create(context.Context, models.CreateTransaction) (string, error)
	Get(context.Context, models.PrimaryKey) (models.Transaction, error)
	GetList(context.Context, models.GetListTransactionsRequest) (models.TransactionsResponse, error)
	Update(context.Context, models.UpdateTransaction) (string, error)
	Delete(context.Context, string) error
}

type IBasketRepo interface {
	Create(context.Context, models.CreateBasket) (string, error)
	Get(context.Context, models.PrimaryKey) (models.Basket, error)
	GetList(context.Context, models.GetBasketsListRequest) (models.BasketsResponse, error)
	Update(context.Context, models.UpdateBasket) (string, error)
	Delete(context.Context, string) error
	UpdateBasketQuantity(context.Context, models.UpdateBasketQuantity) (string, error)
}

type IBranchRepo interface {
	Create(context.Context, models.CreateBranch) (string, error)
	Get(context.Context, models.PrimaryKey) (models.Branch, error)
	GetList(context.Context, models.GetListRequest) (models.BranchsResponse, error)
	Update(context.Context, models.UpdateBranch) (string, error)
	Delete(context.Context, string) error
}

type IProductRepo interface {
	Create(context.Context, models.CreateProduct) (string, error)
	Get(context.Context, models.PrimaryKey) (models.Product, error)
	GetList(context.Context, models.GetListRequest) (models.ProductsResponse, error)
	Update(context.Context, models.UpdateProduct) (string, error)
	Delete(context.Context, string) error
}

type ISaleRepo interface {
	Create(context.Context, models.CreateSale) (string, error)
	Get(context.Context, models.PrimaryKey) (models.Sale, error)
	GetList(context.Context, models.GetListRequest) (models.SalesResponse, error)
	Update(context.Context, models.UpdateSale) (string, error)
	Delete(context.Context, string) error
	UpdateSalePrice(context.Context, models.SaleRequest) (string, error)
}

type IStorageRepo interface {
	Create(context.Context, models.CreateStorage) (string, error)
	Get(context.Context, models.PrimaryKey) (models.Storage, error)
	GetList(context.Context, models.GetListRequest) (models.StoragesResponse, error)
	Update(context.Context, models.UpdateStorage) (string, error)
	Delete(context.Context, string) error
}
