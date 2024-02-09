package postgres

import (
	"bazaar/api/models"
	"bazaar/pkg/check"
	"bazaar/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type staffRepo struct {
	pool *pgxpool.Pool
}

func NewStaffRepo(pool *pgxpool.Pool) storage.IStaffRepo {
	return &staffRepo{
		pool: pool,
	}
}

func (s *staffRepo) Create(ctx context.Context, request models.CreateStaff) (string, error) {

	id := uuid.New()

	query := `insert into staff (
		id, 
		branch_id, 
		tarif_id, 
		type_staff, 
		name, 
		balance, 
		birth_date, 
		age, 
		gender, 
		login, 
		password) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := s.pool.Exec(ctx, query,
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
	)
	if err != nil {
		log.Println("error while inserting staff data", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (s *staffRepo) Get(ctx context.Context, id models.PrimaryKey) (models.Staff, error) {

	var updatedAt = sql.NullTime{}

	staff := models.Staff{}

	query := `select 
	id, 
	branch_id, 
	tarif_id, 
	type_staff, 
	name, 
	birth_date::text, 
	age, 
	gender, 
	login, 
	password, 
	created_at, 
	updated_at from staff where deleted_at is null and id = $1`

	row := s.pool.QueryRow(ctx, query, id.ID)

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
		&updatedAt,
	)

	if err != nil {
		log.Println("error while selecting staff data", err.Error())
		return models.Staff{}, err
	}

	if updatedAt.Valid {
		staff.UpdatedAt = updatedAt.Time
	}

	return staff, nil
}

func (s *staffRepo) GetList(ctx context.Context, request models.GetListRequest) (models.StaffsResponse, error) {
	var (
		updatedAt         = sql.NullTime{}
		staffs            = []models.Staff{}
		count             = 0
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
	)

	countQuery = `select count(1) from staff where deleted_at is null `

	if search != "" {
		countQuery += fmt.Sprintf(`and name ilike '%%%s%%'`, search)
	}
	if err := s.pool.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting staff count", err.Error())
		return models.StaffsResponse{}, err
	}

	query = `select 
	id, 
	branch_id, 
	tarif_id, 
	type_staff, 
	name, 
	birth_date::text, 
	age, 
	gender, 
	login, 
	password, 
	created_at, 
	updated_at from staff where deleted_at is null`

	if search != "" {
		query += fmt.Sprintf(` where name ilike '%%%s%%'`, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := s.pool.Query(ctx, query, request.Limit, offset)
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
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning staff data", err.Error())
			return models.StaffsResponse{}, err
		}

		if updatedAt.Valid {
			staff.UpdatedAt = updatedAt.Time
		}

		staffs = append(staffs, staff)

	}

	return models.StaffsResponse{
		Staffs: staffs,
		Count:  count,
	}, nil
}

func (s *staffRepo) Update(ctx context.Context, request models.UpdateStaff) (string, error) {

	query := `update staff
   set 
   branch_id = $1, 
   tarif_id = $2, 
   type_staff = $3,
   name = $4, 
   birth_date = $5, 
   age = $6, 
   gender = $7, 
   login = $8, 
   password = $9,
   balance = $10, 
   updated_at = $11
   where id = $12
   `
	_, err := s.pool.Exec(ctx, query,
		request.BranchID,
		request.TarifID,
		request.TypeStaff,
		request.Name,
		request.BirthDate,
		check.CalculateAge(request.BirthDate),
		request.Gender,
		request.Login,
		request.Password,
		request.Balance,
		time.Now(),
		request.ID,
	)
	if err != nil {
		log.Println("error while updating staff data...", err.Error())
		return "", err
	}
	return request.ID, nil
}

func (s *staffRepo) Delete(ctx context.Context, id string) error {

	query := `
	update staff
	 set deleted_at = $1
	  where id = $2
	`

	_, err := s.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("error while deleting staff by id", err.Error())
		return err
	}
	return nil
}
