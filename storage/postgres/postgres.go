package postgres

import (
	"bazaar/config"
	"bazaar/storage"
	"context"
	
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database"          //database is needed for migration
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //postgres is used for database
	_ "github.com/golang-migrate/migrate/v4/source/file"       //file is needed for migration url
	

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.Config) (storage.IStorage, error) {
	url := fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		fmt.Println("error while parsing config", err.Error())
		return nil, err
	}

	poolConfig.MaxConns = 100

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		fmt.Println("error while connecting to db", err.Error())
		return nil, err
	}

	return Store{
		pool: pool,
	}, nil

}

func (s Store) CloseDB() {
	s.pool.Close()
}

func (s Store) Category() storage.ICategoryRepo {
	return NewCategoryRepo(s.pool)
}

func (s Store) Staff() storage.IStaffRepo {
	return NewStaffRepo(s.pool)
}

func (s Store) StorageTransaction() storage.IStorageTransactionRepo {
	return NewStorageTransactionRepo(s.pool)
}

func (s Store) Tarif() storage.ITarifRepo {
	return NewTarifRepo(s.pool)
}

func (s Store) Transaction() storage.ITransactionRepo {
	return NewTransactionRepo(s.pool)
}

func (s Store) Basket() storage.IBasketRepo {
	return NewBasketRepo(s.pool)
}

func (s Store) Branch() storage.IBranchRepo {
	return NewBranchRepo(s.pool)
}

func (s Store) Product() storage.IProductRepo {
	return NewProductRepo(s.pool)
}

func (s Store) Sale() storage.ISaleRepo {
	return NewSaleRepo(s.pool)
}

func (s Store) Storage() storage.IStorageRepo {
	return NewStorageRepo(s.pool)
}
