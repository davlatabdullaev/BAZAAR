package postgres

import (
	"bazaar/api/models"
	"bazaar/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tarifRepo struct {
	pool *pgxpool.Pool
}

func NewTarifRepo(pool *pgxpool.Pool) storage.ITarifRepo {
	return &tarifRepo{
		pool: pool,
	}
}

func (t *tarifRepo) Create(ctx context.Context, request models.CreateTarif) (string, error) {

	id := uuid.New()

	query := `insert into tarif (id, name, tarif_type, amount_for_cash,
		amount_for_card) 
	values 
	($1, $2, $3, $4, $5)`

	_, err := t.pool.Exec(ctx, query,
		id,
		request.Name,
		request.TarifType,
		request.AmountForCash,
		request.AmountForCard,
	)
	if err != nil {
		log.Println("error while inserting tarif data", err.Error())
		return "", err
	}

	return id.String(), nil

}

func (t *tarifRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Tarif, error) {

	var updatedAt = sql.NullTime{}

	tarif := models.Tarif{}

	query := `select 
	id, 
	name, 
	tarif_type, 
	amount_for_cash,
	amount_for_card, 
	created_at, 
	updated_at from tarif
	 where deleted_at is null and id = $1`

	row := t.pool.QueryRow(ctx, query, id.ID)

	err := row.Scan(
		&tarif.ID,
		&tarif.Name,
		&tarif.TarifType,
		&tarif.AmountForCash,
		&tarif.AmountForCard,
		&tarif.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		log.Println("error while selecting tarif data", err.Error())
		return models.Tarif{}, err
	}

	if updatedAt.Valid {
		tarif.UpdatedAt = updatedAt.Time
	}

	return tarif, nil

}

func (t *tarifRepo) GetList(ctx context.Context, request models.GetListRequest) (models.TarifsResponse, error) {

	var (
		updatedAt         = sql.NullTime{}
		tarifs            = []models.Tarif{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from tarif where deleted_at is null `

	if search != "" {
		countQuery += fmt.Sprintf(`and name ilike '%%%s%%'`, search)
	}
	if err := t.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting tarif count", err.Error())
		return models.TarifsResponse{}, err
	}

	query = `select id, name, tarif_type, 
	amount_for_cash, amount_for_card,
	created_at, updated_at
	from tarif 
	where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` where name ilike '%%%s%%'`, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := t.pool.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting tarif", err.Error())
		return models.TarifsResponse{}, err
	}

	for rows.Next() {
		tarif := models.Tarif{}
		if err = rows.Scan(
			&tarif.ID,
			&tarif.Name,
			&tarif.TarifType,
			&tarif.AmountForCash,
			&tarif.AmountForCard,
			&tarif.CreatedAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning tarif data", err.Error())
			return models.TarifsResponse{}, err
		}

		if updatedAt.Valid {
			tarif.UpdatedAt = updatedAt.Time
		}

		tarifs = append(tarifs, tarif)

	}

	return models.TarifsResponse{
		Tarifs: tarifs,
		Count:  count,
	}, nil
}

func (t *tarifRepo) Update(ctx context.Context, request models.UpdateTarif) (string, error) {

	query := `update tarif
   set name = $1, tarif_type = $2, amount_for_cash = $3,
   amount_for_card = $4, updated_at = $5
   where id = $6
   `
	_, err := t.pool.Exec(ctx, query,
		request.Name,
		request.TarifType,
		request.AmountForCash,
		request.AmountForCard,
		time.Now(),
		request.ID,
	)
	if err != nil {
		log.Println("error while updating tarif data...", err.Error())
		return "", err
	}
	return request.ID, nil
}

func (t *tarifRepo) Delete(ctx context.Context, id string) error {

	query := `
	update tarif
	 set deleted_at = $1
	  where id = $2
	`

	_, err := t.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting tarif by id", err.Error())
		return err
	}

	return nil
}
