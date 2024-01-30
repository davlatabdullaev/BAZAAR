package postgres

import (
	"bazaar/config"
	"bazaar/storage"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func New(cfg config.Config) (storage.IStorage, error) {
	url := fmt.Sprintf(`host = %s port = %s user = %s password = %s database = %s sslmode=disable`, 
     cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)

	 db, err := sql.Open("postgres", url)
	 if err != nil {
		return Store{}, err
	 }

	 return Store{
		db: db,
	 }, nil

}

func (s Store) CloseDB() {
	s.db.Close()
}

func (s Store) Category() storage.ICategoryRepo {
	return NewCategoryRepo(s.db)
}

func (s Store) Staff() storage.IStaffRepo {
	return  NewStaffRepo(s.db)
}

func (s Store) StorageTransaction() storage.IStorageTransactionRepo {
	return NewStorageTransactionRepo(s.db)
}

func (s Store) Tarif() storage.ITarifRepo{
	return NewTarifRepo(s.db)
}

func (s Store) Transaction() storage.ITransactionRepo {
   return NewTransactionRepo(s.db)
}