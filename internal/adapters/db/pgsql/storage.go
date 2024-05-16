package pgsql

import (
	"context"
	"goshell/internal/entities"
	"goshell/internal/pgclient"
	"log"
	"os"
	"os/exec"
	"time"

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

// type ResultStorage struct
type ResultStorage struct {
	dbpool *pgxpool.Pool
}

// func NewResultStorage() *ResultStorage
func NewResultStorage() *ResultStorage {
	return &ResultStorage{dbpool: pgclient.WDB}
}

// func (est *CommandStorage) CommExec(ctx context.Context, id int) error
func (est *CommandStorage) CommExec(ctx context.Context, id int) error {
	g := entities.Command{}

	err := est.dbpool.QueryRow(ctx, "SELECT * from public.commands where id=$1;", id).Scan(&g.Id, &g.CommandText, &g.ScriptText)

	if err != nil {
		log.Println(err.Error(), "exec_one")
	}

	_ = os.MkdirAll("/test", 0777)

	f, err := os.CreateTemp("/test", "*.sh")
	if err != nil {
		log.Println(err.Error())
	}
	defer f.Close()

	if err := os.WriteFile(f.Name(), []byte(g.ScriptText), 0777); err != nil {
		log.Println(err.Error())
	}

	stout := ""
	cmd := f.Name()

	out, err := exec.Command("bash", cmd).Output()

	if err != nil {
		log.Println(err.Error())
	}

	stout = string(out)

	rid := 0
	t := time.Now()
	err = est.dbpool.QueryRow(ctx, "insert into results (id, id_command, output, time) values (default, $1, $2, $3) returning id;", g.Id, stout, t.Format("2006-01-02 15:04:05")).Scan(&rid)

	if err != nil {
		log.Println("Failed execute command add!")
		return err
	}

	return nil
}

// func (est *CommandStorage) CommSave(ctx context.Context, bs []byte) (int, error)
func (est *CommandStorage) CommSave(ctx context.Context, bs []byte) (int, error) {
	_ = os.MkdirAll("/test", 0777)

	f, err := os.CreateTemp("/test", "*.sh")
	if err != nil {
		log.Println(err.Error())
	}
	defer f.Close()

	if err := os.WriteFile(f.Name(), bs, 0777); err != nil {
		log.Println(err.Error())
	}

	cid := 0
	err = est.dbpool.QueryRow(ctx, "insert into commands (id, command_text, script_text) values (default, $1, $2) returning id;", f.Name(), string(bs)).Scan(&cid)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return cid, nil
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

// func (est *ResultStorage) GetList(ctx context.Context) (entities.Result_count, error)
func (est *ResultStorage) GetList(ctx context.Context) (entities.Result_count, error) {
	gs := entities.Result{}

	gsc := 0
	err := est.dbpool.QueryRow(ctx, "SELECT count(*) from public.results;").Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "results_count")
		return entities.Result_count{Values: []entities.Result{}, Count: 0}, err
	}

	out_arr := make([]entities.Result, 0, gsc)

	rows, err := est.dbpool.Query(ctx, "SELECT id, id_command, output, time::text as ts from public.results;")
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
