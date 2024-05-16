package pgsql

import (
	"context"
	"goshell/internal/entities"
	"goshell/internal/pgclient"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//type CommandStorage struct
type CommandStorage struct {
	dbpool *pgxpool.Pool
}

//func NewCommandStorage() *CommandStorage
func NewCommandStorage() *CommandStorage {
	return &CommandStorage{dbpool: pgclient.WDB}
}

// func (est *CommandStorage) GetList(ctx context.Context) (entities.Command_count, error)
func (est *CommandStorage) GetList(ctx context.Context) (entities.Command_count, error) {
	gs := entities.Command{}

	gsc := 0
	err := est.dbpool.QueryRow(ctx, "SELECT count(*) from public.commands;").Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "commands_count")
		return entities.Command_count{Values: []entities.Command{}, Count: 0}, err
	}

	out_arr := make([]entities.Command, 0, gsc)

	rows, err := est.dbpool.Query(ctx, "SELECT * from public.commands;")
	if err != nil {
		log.Println(err.Error(), "commands_list")
		return entities.Command_count{Values: []entities.Command{}, Count: 0}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.CommandText, &gs.ScriptText)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_arr_count := entities.Command_count{Values: out_arr, Count: gsc}

	return out_arr_count, nil
}

// func (est *CommandStorage) GetOne(ctx context.Context, id int) (entities.Command_count, error)
func (est *CommandStorage) GetOne(ctx context.Context, id int) (entities.Command_count, error) {
	out_arr := []entities.Command{}
	g := entities.Command{}

	err := est.dbpool.QueryRow(ctx, "SELECT * from public.commands where id=$1;", id).Scan(&g.Id, &g.CommandText, &g.ScriptText)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute commands_one: ", err)
		return entities.Command_count{Values: []entities.Command{}, Count: 0}, err
	}

	out_arr = append(out_arr, g)

	out_count := entities.Command_count{Values: out_arr, Count: 1}
	return out_count, nil
}
