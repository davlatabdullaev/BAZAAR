package postgres

import (
	"bazaar/api/models"
	"bazaar/storage"
	"database/sql"
)

type branchRepo struct {
	db *sql.DB
}

func NewBranchRepo(db *sql.DB) storage.IBranchRepo {
	return branchRepo {
		db: db,
	}
}

func (b branchRepo) Create(branch models.CreateBranch) (string, error) {

	return "", nil
}

func (b branchRepo) Get(id models.PrimaryKey) (models.Branch, error) {

	return models.Branch{}, nil
}

func (b branchRepo) GetList(request models.GetListRequest) (models.BranchsResponse, error) {

	return models.BranchsResponse{}, nil
}

func (b branchRepo) Update(request models.UpdateBranch) (string, error) {

	return "", nil
}

func (b branchRepo) Delete(id string) error {

	return nil
}
