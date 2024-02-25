package postgres

import (
	"bazaar/config"
	"bazaar/pkg/logger"
	"bazaar/storage"
	"context"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database"          //database is needed for migration
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //postgres is used for database
	_ "github.com/golang-migrate/migrate/v4/source/file"       //file is needed for migration url

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type Store struct {
	pool *pgxpool.Pool
	log  logger.ILogger
}

func New(ctx context.Context, cfg config.Config, log logger.ILogger) (storage.IStorage, error) {
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
		log.Error("error while parsing config", logger.Error(err))
		return nil, err
	}

	poolConfig.MaxConns = 100

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error("error while connecting to db", logger.Error(err))
		return nil, err
	}

	// migration
	m, err := migrate.New("file://migrations/postgres/", url)
	if err != nil {
		log.Error("error while migrating", logger.Error(err))
		return nil, err
	}

	if err = m.Up(); err != nil {
		fmt.Println("here up")
		if !strings.Contains(err.Error(), "no change") {
			fmt.Println("in !strings")
			version, dirty, err := m.Version()
			if err != nil {
				log.Error("error in checking version and dirty", logger.Error(err))
				return nil, err
			}

			if dirty {
				version--
				if err = m.Force(int(version)); err != nil {
					log.Error("error in making force", logger.Error(err))
					return nil, err
				}
			}
		}
	}

	return Store{
		pool: pool,
		log:  log,
	}, nil

}

func (s Store) CloseDB() {
	s.pool.Close()
}

func (s Store) Category() storage.ICategoryRepo {
	return NewCategoryRepo(s.pool, s.log)
}

func (s Store) Staff() storage.IStaffRepo {
	return NewStaffRepo(s.pool, s.log)
}

func (s Store) StorageTransaction() storage.IStorageTransactionRepo {
	return NewStorageTransactionRepo(s.pool, s.log)
}

func (s Store) Tarif() storage.ITarifRepo {
	return NewTarifRepo(s.pool, s.log)
}

func (s Store) Transaction() storage.ITransactionRepo {
	return NewTransactionRepo(s.pool, s.log)
}

func (s Store) Basket() storage.IBasketRepo {
	return NewBasketRepo(s.pool, s.log)
}

func (s Store) Branch() storage.IBranchRepo {
	return NewBranchRepo(s.pool, s.log)
}

func (s Store) Product() storage.IProductRepo {
	return NewProductRepo(s.pool, s.log)
}

func (s Store) Sale() storage.ISaleRepo {
	return NewSaleRepo(s.pool, s.log)
}

func (s Store) Storage() storage.IStorageRepo {
	return NewStorageRepo(s.pool, s.log)
}

func (s Store) Income() storage.IIncomeRepo {
	return NewIncomeRepo(s.pool, s.log)
}

func (s Store) IncomeProduct() storage.IIncomeProductRepo {
	return NewIncomeProductRepo(s.pool, s.log)
}
