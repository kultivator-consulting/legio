package database

import (
	"context"
	"cortex_api/database/db_gen"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"log"
	"os"
	"strconv"
)

const DefaultDatabasePort = 5432

type Model struct {
	DbPool *pgxpool.Pool
}

type Interface interface {
	Open() (*pgxpool.Pool, *db_gen.Queries, context.Context, error)
	Close() error
}

var Module = fx.Options(fx.Provide(ApiDatabase))

func ApiDatabase() *Model {
	return &Model{
		DbPool: nil,
	}
}

func (db *Model) Open() (*pgxpool.Pool, *db_gen.Queries, context.Context, error) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = DefaultDatabasePort
	}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	psqlConnection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

	log.Printf("Opening database connection.\n")

	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, psqlConnection)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, nil, ctx, err
	}

	db.DbPool = dbPool
	queries := db_gen.New(dbPool)

	return dbPool, queries, ctx, nil
}

func (db *Model) Close() error {
	log.Printf("Closing database connection.\n")
	if db.DbPool == nil {
		return nil
	}
	return db.Close()
}
