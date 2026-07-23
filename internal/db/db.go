package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect(database_url string) error {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, database_url)
	if err != nil {
		log.Println("Error creating pool")
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Println("Error connecting database")
		return err
	}
	log.Println("Connected succesfully")
	Pool = pool
	return nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
		log.Println("Closed succesfully")
	}
}
