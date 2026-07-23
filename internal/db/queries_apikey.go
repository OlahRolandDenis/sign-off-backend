package db

import (
	"Desktop/signoff/internal/models"
	"context"
	"fmt"
	"log"
)

func CreateAPIKey(ctx context.Context, agencyID int64, name, key string) (int64, error) {
	var query = `INSERT INTO api_keys (agency_id, key, name) VALUES ($1, $2, $3) RETURNING id`
	row := Pool.QueryRow(ctx, query, agencyID, key, name)

	var id int64
	res := row.Scan(&id)
	if res != nil {
		log.Printf("Error getting id: %v", res)
		return 0, res
	}
	return id, nil
}

func GetAPIKeys(ctx context.Context, agencyID int64) ([]models.APIKey, error) {
	var query = `SELECT id, agency_id, key, name, created_at, last_used FROM api_keys WHERE agency_id=$1 ORDER BY created_at DESC`
	rows, err := Pool.Query(ctx, query, agencyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []models.APIKey
	for rows.Next() {
		var k models.APIKey
		err = rows.Scan(&k.ID, &k.AgencyID, &k.Key, &k.Name, &k.CreatedAt, &k.LastUsed)
		if err != nil {
			continue
		}
		keys = append(keys, k)
	}
	return keys, nil

}

func DeleteAPIKey(ctx context.Context, id, agencyID int64) error {
	var query = `DELETE FROM api_keys WHERE id=$1 and agency_id=$2`
	result, err := Pool.Exec(ctx, query, id)

	if err != nil {
		log.Printf("Error deleting id: %v", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("API key not found or not owned by you")
	}
	return nil
}

func ValidateAPIKey(ctx context.Context, key string) (int64, error) {
	var query = `SELECT agency_id FROM api_keys WHERE key=$1`
	row := Pool.QueryRow(ctx, query, key)
	var id int64

	res := row.Scan(&id)
	if res != nil {
		log.Printf("Error getting id: %v", res)
		return 0, res
	}

	return id, nil
}
