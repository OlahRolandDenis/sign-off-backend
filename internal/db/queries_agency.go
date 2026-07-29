package db

import (
	"Desktop/signoff/internal/models"
	"context"
	"log"
	"time"
)

func CreateAgency(ctx context.Context, email, name, password_hash string) (int64, error) {
	var query = `INSERT INTO agencies (email, name, password_hash) VALUES ($1, $2, $3) RETURNING id`
	row := Pool.QueryRow(ctx, query, email, name, password_hash)

	var id int64
	res := row.Scan(&id)
	if res != nil {
		log.Printf("Error getting id: %v", res)
		return 0, res
	}

	return id, nil
}

func GetAgencyByEmail(ctx context.Context, email string) (*models.Agency, error) {
	var query = `SELECT * FROM agencies WHERE email=$1`
	row := Pool.QueryRow(ctx, query, email)

	var id int64
	var name string
	var emailul string
	var password_hash string
	var plan string
	var created_at time.Time

	err := row.Scan(&id, &emailul, &name, &password_hash, &plan, &created_at)
	if err != nil {
		log.Printf("Error finding agency id: %v", err)
		return nil, err
	}

	var ma models.Agency
	ma.CreatedAt = created_at
	ma.Email = emailul
	ma.ID = id
	ma.Name = name
	ma.PasswordHash = password_hash
	ma.Plan = plan

	return &ma, nil
}

func GetAgencyByID(ctx context.Context, id int64) (*models.Agency, error) {
	var query = `SELECT * FROM agencies WHERE id=$1`
	row := Pool.QueryRow(ctx, query, id)

	var idul int64
	var name string
	var emailul string
	var password_hash string
	var plan string
	var created_at time.Time

	err := row.Scan(&idul, &emailul, &name, &password_hash, &plan, &created_at)
	if err != nil {
		log.Printf("Error finding agency id: %v", err)
		return nil, err
	}
	var ma models.Agency
	ma.CreatedAt = created_at
	ma.Email = emailul
	ma.ID = idul
	ma.Name = name
	ma.PasswordHash = password_hash
	ma.Plan = plan

	return &ma, nil
}

func DeleteAgency(ctx context.Context, id int64) error {
	var query = `DELETE * FROM agencies WHERE id=$1`
	_, err := Pool.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Error deleting agency by id: %v", err)
		return err
	}
	return nil
}

func CountAPIKeys(ctx context.Context, agencyID int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM api_keys WHERE agency_id = $1`
	err := Pool.QueryRow(ctx, query, agencyID).Scan(&count)
	return count, err
}

func CountApprovalsInLast30Days(ctx context.Context, agencyID int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM approvals WHERE agency_id = $1 AND created_at > NOW() - INTERVAL '30 days'`
	err := Pool.QueryRow(ctx, query, agencyID).Scan(&count)
	return count, err
}

func GetPlan(ctx context.Context, agencyID int64) (string, error) {
	var plan string
	query := `SELECT plan FROM agencies WHERE id = $1`
	err := Pool.QueryRow(ctx, query, agencyID).Scan(&plan)
	return plan, err
}
