package services

import (
	"context"
	"goshell/internal/entities"
	"goshell/internal/pgclient"
	"log"
)

// func ResultGetList(ctx context.Context) (entities.Result_count, error) - вывод списка результатов
func ResultGetList(ctx context.Context) (entities.Result_count, error) {
	dbpool := pgclient.WDB
	gs := entities.Result{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT count(*) from public.results;").Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "results_count")
		return entities.Result_count{Values: []entities.Result{}, Count: 0}, err
	}

	out_arr := make([]entities.Result, 0, gsc)

	rows, err := dbpool.Query(ctx, "SELECT id, id_command, output, time::text as ts from public.results;")
	if err != nil {
		log.Println(err.Error(), "results_list")
		return entities.Result_count{Values: []entities.Result{}, Count: 0}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.IdC, &gs.Output, &gs.TS)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_arr_count := entities.Result_count{Values: out_arr, Count: gsc}

	return out_arr_count, nil
}
