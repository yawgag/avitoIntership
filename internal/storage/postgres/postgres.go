package postgres

import (
	"context"
	"log"
	"orderPickupPoint/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDb() (*pgxpool.Pool, error) {
	serverConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbConfig, err := pgxpool.ParseConfig(serverConfig.DbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbConfig.MaxConns = 20
	dbConfig.MinConns = 1
	dbConfig.MaxConnLifetime = time.Hour

	DbConnPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = DbConnPool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return DbConnPool, err
}

func CloseDB(DbConnPool *pgxpool.Pool) {
	if DbConnPool != nil {
		DbConnPool.Close()
	}
}
