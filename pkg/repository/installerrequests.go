package repository

import (
	"database/sql"
)

type InstallerRequest struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	State       string `json:"state"`
	City        string `json:"city"`
	PlanId      string `json:"plan_id"`
	InstallerId string `json:"installer_id"`
}

type InstallerRequestsRepository struct {
	DB *sql.DB
}

func NewInstallerRequestsRepository(db *sql.DB) *InstallerRequestsRepository {
	return &InstallerRequestsRepository{
		DB: db,
	}
}

func (ir *InstallerRequestsRepository) CreateRequest(request InstallerRequest) error {
	stmt := `INSERT INTO requests (email, name, state, city, plan_id, installer_id)
    VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := ir.DB.Exec(stmt, request.Email, request.Name, request.State, request.City, request.PlanId, request.InstallerId)
	if err != nil {
		return err
	}

	return nil
}

func (ir *InstallerRequestsRepository) GetRequestsByPlanAndInstallerId(planId, installerId int) ([]InstallerRequest, error) {
	stmt := `SELECT * from requests WHERE plan_id = $1 AND installer_id = $2`

	rows, err := ir.DB.Query(stmt, planId, installerId)
	if err != nil {
		return []InstallerRequest{}, err
	}
	defer rows.Close()

	var requests []InstallerRequest
	for rows.Next() {
		var request InstallerRequest
		err = rows.Scan(&request.Email, &request.Name, &request.State, &request.City, &request.PlanId, &request.InstallerId)
		if err != nil {
			panic(err)
		}
		requests = append(requests, request)
	}

	return requests, nil
}
