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

	// TODO: migrations

	ct_sql := `CREATE TABLE IF NOT EXISTS public.commands (
		id int NOT NULL GENERATED ALWAYS AS IDENTITY,
		command_text text NOT NULL,
		result_text text NOT NULL,
		CONSTRAINT command_pk PRIMARY KEY (id)
	);`

	_, err = WDB.Exec(ctx, ct_sql)
	if err != nil {
		log.Printf("Error %s when creating commands table", err)
		return nil, err
	}

	return WDB, nil
}
