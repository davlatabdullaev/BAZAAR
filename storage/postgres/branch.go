package postgres

import (
	"bazaar/api/models"
	"bazaar/pkg/check"
	"bazaar/storage"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type branchRepo struct {
	db *sql.DB
}

func NewBranchRepo(db *sql.DB) storage.IBranchRepo {
	return branchRepo{
		db: db,
	}
}

func (b branchRepo) Create(branch models.CreateBranch) (string, error) {

	id := uuid.New()

	query := `insert into category (id, name, address, updated_at) values ($1, $2, $3, $4)`

	res, err := b.db.Exec(query,
		id,
		branch.Name,
		branch.Address,
		check.TimeNow(),
	)
	if err != nil {
		log.Println("error while inserting branch", err.Error())
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

func (b branchRepo) Get(id models.PrimaryKey) (models.Branch, error) {

	branch := models.Branch{}

	row := b.db.QueryRow(`select id, name, address, created_at, updated_at from branch where deleted_at is null and id = $1`, id)

	err := row.Scan(
		&branch.ID,
		&branch.Name,
		&branch.Address,
		&branch.CreatedAt,
		&branch.UpdatedAt,
	)

	if err != nil {
		log.Println("error while selecting branch", err.Error())
		return models.Branch{}, err
	}

	return branch, nil

}

func (b branchRepo) GetList(request models.GetListRequest) (models.BranchsResponse, error) {

	var (
		branchs           = []models.Branch{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from branch `

	if search != "" {
		countQuery += fmt.Sprintf(`where name ilike '%%%s%%'`, search)
	}
	if err := b.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
		return models.BranchsResponse{}, err
	}

	query = `select id, name, address, created_at, updated_at from branch where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` where name ilike '%%%s%%'`, search)
	}

	query += `LIMIT $1 OFFSET $2`
	rows, err := b.db.Query(query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting branch", err.Error())
		return models.BranchsResponse{}, err
	}

	for rows.Next() {
		branch := models.Branch{}
		if err = rows.Scan(
			&branch.ID,
			&branch.Name,
			&branch.Address,
			&branch.CreatedAt,
			&branch.UpdatedAt,
		); err != nil {
			fmt.Println("error is while scanning branch data", err.Error())
			return models.BranchsResponse{}, err
		}

		branchs = append(branchs, branch)

	}

	return models.BranchsResponse{
		Branchs: branchs,
		Count:   count,
	}, nil
}

func (b branchRepo) Update(request models.UpdateBranch) (string, error) {

	query := `update branch
   set name = $1,
    address = $2, 
	updated_at = $3
   where id = $4  
   `

	res, err := b.db.Exec(query,
		request.Name,
		request.Address,
		request.UpdatedAt,
		request.ID)
	if err != nil {
		log.Println("error while updating branch data...", err.Error())
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

func (b branchRepo) Delete(id string) error {

	query := `
	update branch
	 set deleted_at = $1
	  where id = $2
	`

	res, err := b.db.Exec(query, check.TimeNow(), id)
	if err != nil {
		log.Println("error while deleting branch by id", err.Error())
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
