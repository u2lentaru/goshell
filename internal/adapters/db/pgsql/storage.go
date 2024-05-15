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

//func NewCommandStorage(db *pgxpool.Pool) *CommandStorage
func NewCommandStorage() *CommandStorage {
	return &CommandStorage{dbpool: pgclient.WDB}
}

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
