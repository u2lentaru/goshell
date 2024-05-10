package pgclient

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var WDB *pgxpool.Pool

func GetDb(ctx context.Context, url string) (*pgxpool.Pool, error) {
	if WDB != nil {
		return WDB, nil
	}

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(err)
	}

	cfg.MaxConns = 8
	cfg.MinConns = 1

	WDB, err = pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	// defer dbpool.Close()

	rows, err := WDB.Query(ctx, "SELECT version();")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		v := ""
		err = rows.Scan(&v)

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		log.Println("version:", v)
	}

	return WDB, nil
}
