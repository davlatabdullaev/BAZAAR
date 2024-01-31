package postgres

import (
	"bazaar/api/models"
	"bazaar/pkg/check"
	"bazaar/storage"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type tarifRepo struct {
	db *sql.DB
}

func NewTarifRepo(db *sql.DB) storage.ITarifRepo {
	return tarifRepo{
		db: db,
	}
}

func (t tarifRepo) Create(request models.CreateTarif) (string, error) {

	id := uuid.New()

	updatedAt := time.Now()

	query := `insert into tarif (id, name, tarif_type, amount_for_cash,
		amount_for_card, updated_at) 
	values 
	($1, $2, $3, $4, $5, $6)`

	res, err := t.db.Exec(query,
		id,
		request.Name,
		request.TarifType,
		request.AmountForCash,
		request.AmountForCard,
		updatedAt,
	)
	if err != nil {
		log.Println("error while inserting tarif data", err.Error())
		return "", err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("error getting rows affected", err.Error())
		return "", err
	}

	if rowsAffected == 0 {
		log.Println("no rows affected during insert")
		return "", errors.New("no rows affected during insert")
	}

	return id.String(), nil

}

func (t tarifRepo) Get(id models.PrimaryKey) (models.Tarif, error) {

	tarif := models.Tarif{}

	query := `select id, name, tarif_type, amount_for_cash,
	amount_for_card, 
	 created_at, updated_at from tarif
	 where deleted_at = null and id = $1`

	row := t.db.QueryRow(query, id)

	err := row.Scan(
		&tarif.ID,
		&tarif.Name,
		&tarif.TarifType,
		&tarif.AmountForCash,
		&tarif.AmountForCard,
		&tarif.CreatedAt,
		&tarif.UpdatedAt,
	)

	if err != nil {
		log.Println("error while selecting tarif data", err.Error())
		return models.Tarif{}, err
	}

	return tarif, nil

}

func (t tarifRepo) GetList(request models.GetListRequest) (models.TarifsResponse, error) {

	var (
		tarifs            = []models.Tarif{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from tarif `

	if search != "" {
		countQuery += fmt.Sprintf(`where name ilike '%%%s%%'`, search)
	}
	if err := t.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting tarif count", err.Error())
		return models.TarifsResponse{}, err
	}

	query = `select id, name, tarif_type, 
	amount_for_cash, amount_for_card,
	created_at, updated_at
	from tarif 
	where deleted_at = null`

	if search != "" {
		query += fmt.Sprintf(` where name ilike '%%%s%%'`, search)
	}

	query += `LIMIT $1 OFFSET $2`
	rows, err := t.db.Query(query, request.Limit, offset)
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
			&tarif.UpdatedAt,
		); err != nil {
			fmt.Println("error is while scanning tarif data", err.Error())
			return models.TarifsResponse{}, err
		}

		tarifs = append(tarifs, tarif)

	}

	return models.TarifsResponse{
		Tarifs: tarifs,
		Count:  count,
	}, nil
}

func (t tarifRepo) Update(request models.UpdateTarif) (string, error) {

	query := `update tarif
   set name = $1, tarif_type = $2, amount_for_cash = $3,
   amount_for_card = $4, updated_at = $5
   where id = $6
   `
	res, err := t.db.Exec(query,
		request.Name,
		request.TarifType,
		request.AmountForCash,
		request.AmountForCard,
		check.TimeNow(),
		request.ID,
	)
	if err != nil {
		log.Println("error while updating tarif data...", err.Error())
		return "", err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("error getting rows affected", err.Error())
		return "", err
	}

	if rowsAffected == 0 {
		log.Println("no rows affected during insert")
		return "", errors.New("no rows affected during insert")
	}

	return "", nil
}

func (t tarifRepo) Delete(id string) error {

	query := `
	update tarif
	 set deleted_at = $1
	  where id = $2
	`

	res, err := t.db.Exec(query, check.TimeNow(), id)
	if err != nil {
		log.Println("error while deleting tarif by id", err.Error())
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("error getting rows affected", err.Error())
		return err
	}

	if rowsAffected == 0 {
		log.Println("no rows affected during insert")
		return errors.New("no rows affected during insert")
	}

	return nil
}
