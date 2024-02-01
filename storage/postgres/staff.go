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

type staffRepo struct {
	db *sql.DB
}

func NewStaffRepo(db *sql.DB) storage.IStaffRepo {
	return staffRepo{
		db: db,
	}
}

func (s staffRepo) Create(request models.CreateStaff) (string, error) {

	id := uuid.New()

	updatedAt := time.Now()

	query := `insert into staff (id, branch_id, tarif_id, type_staff, name, balance, birth_date, age, gender, login, password, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	res, err := s.db.Exec(query,
		id,
		request.BranchID,
		request.TarifID,
		request.TypeStaff,
		request.Name,
		request.Balance,
		request.BirthDate,
		check.CalculateAge(request.BirthDate),
		request.Gender,
		request.Login,
		request.Password,
		updatedAt,
	)
	if err != nil {
		log.Println("error while inserting staff data", err.Error())
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

func (s staffRepo) Get(id models.PrimaryKey) (models.Staff, error) {

	staff := models.Staff{}

	query := `select id, branch_id, tarif_id, type_staff, name, birth_date, age, gender, login, password, created_at, updated_at from staff where deleted_at is null and id = $1`

	row := s.db.QueryRow(query, id.ID)

	err := row.Scan(
		&staff.ID,
		&staff.BranchID,
		&staff.TarifID,
		&staff.TypeStaff,
		&staff.Name,
		&staff.BirthDate,
		&staff.Age,
		&staff.Gender,
		&staff.Login,
		&staff.Password,
		&staff.CreatedAt,
		&staff.UpdatedAt,
	)

	if err != nil {
		log.Println("error while selecting staff data", err.Error())
		return models.Staff{}, err
	}

	return staff, nil
}

func (s staffRepo) GetList(request models.GetListRequest) (models.StaffsResponse, error) {
	var (
		staffs            = []models.Staff{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from staff `

	if search != "" {
		countQuery += fmt.Sprintf(`where name ilike '%%%s%%'`, search)
	}
	if err := s.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting staff count", err.Error())
		return models.StaffsResponse{}, err
	}

	query = `select id, branch_id, tarif_id, type_staff, name, birth_date, age, gender, login, password, created_at, updated_at from staff where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` where name ilike '%%%s%%'`, search)
	}

	query += `LIMIT $1 OFFSET $2`
	rows, err := s.db.Query(query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting staff", err.Error())
		return models.StaffsResponse{}, err
	}

	for rows.Next() {
		staff := models.Staff{}
		if err = rows.Scan(
			&staff.ID,
			&staff.BranchID,
			&staff.TarifID,
			&staff.TypeStaff,
			&staff.Name,
			&staff.BirthDate,
			&staff.Age,
			&staff.Gender,
			&staff.Login,
			&staff.Password,
			&staff.CreatedAt,
			&staff.UpdatedAt,
		); err != nil {
			fmt.Println("error is while scanning staff data", err.Error())
			return models.StaffsResponse{}, err
		}

		staffs = append(staffs, staff)

	}

	return models.StaffsResponse{
		Staffs: staffs,
		Count:  count,
	}, nil
}

func (s staffRepo) Update(request models.UpdateStaff) (string, error) {

	query := `update staff
   set branch_id = $1, tarif_id = $2, type_staff = $3,
   name = $4, birth_date = $5, age = $6, gender = $7, login = $8, password = $9, updated_at = $10
   where id = $11
   `
	res, err := s.db.Exec(query,
		request.BranchID,
		request.TarifID,
		request.TypeStaff,
		request.Name,
		request.BirthDate,
		check.CalculateAge(request.BirthDate),
		request.Gender,
		request, request.Login,
		request.Password,
		request.UpdatedAt,
		request.ID,
	)
	if err != nil {
		log.Println("error while updating category data...", err.Error())
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

	return request.ID, nil
}

func (s staffRepo) Delete(id string) error {

	query := `
	update staff
	 set deleted_at = $1
	  where id = $2
	`

	res, err := s.db.Exec(query, check.TimeNow(), id)
	if err != nil {
		log.Println("error while deleting staff by id", err.Error())
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
