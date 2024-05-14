package services

import (
	"context"
	"goshell/internal/entities"
	"goshell/internal/pgclient"
	"log"
	"os"
	"os/exec"
	"time"
)

// func CommExec(id int) - выполняет скрипт и сохраняет результат
func CommExec(id int) error {
	ctx := context.Background()
	dbpool := pgclient.WDB
	g := entities.Command{}

	err := dbpool.QueryRow(ctx, "SELECT * from public.commands where id=$1;", id).Scan(&g.Id, &g.CommandText, &g.ScriptText)

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
	err = dbpool.QueryRow(ctx, "insert into results (id, id_command, output, time) values (default, $1, $2, $3) returning id;", g.Id, stout, t.Format("2006-01-02 15:04:05")).Scan(&rid)

	if err != nil {
		log.Println("Failed execute command add!")
		return err
	}

	return nil
}

// func CommSave(bs []byte) (int, error) - сохраняет скрипт в базу
func CommSave(bs []byte) (int, error) {
	ctx := context.Background()
	dbpool := pgclient.WDB

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
	err = dbpool.QueryRow(ctx, "insert into commands (id, command_text, script_text) values (default, $1, $2) returning id;", f.Name(), string(bs)).Scan(&cid)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return cid, nil
}

// func CommGetList(ctx context.Context) (entities.Command_count, error)
func CommGetList(ctx context.Context) (entities.Command_count, error) {
	dbpool := pgclient.WDB
	gs := entities.Command{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT count(*) from public.commands;").Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "commands_count")
		return entities.Command_count{Values: []entities.Command{}, Count: 0}, err
	}

	out_arr := make([]entities.Command, 0, gsc)

	rows, err := dbpool.Query(ctx, "SELECT * from public.commands;")
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

// func CommGetOne(ctx context.Context, id int) (entities.Command_count, error)
func CommGetOne(ctx context.Context, id int) (entities.Command_count, error) {
	dbpool := pgclient.WDB
	out_arr := []entities.Command{}
	g := entities.Command{}

	err := dbpool.QueryRow(ctx, "SELECT * from public.commands where id=$1;", id).Scan(&g.Id, &g.CommandText, &g.ScriptText)

	if err != nil {
		log.Println(err.Error(), "commands_one")
		return entities.Command_count{Values: []entities.Command{}, Count: 0}, err
	}

	out_arr = append(out_arr, g)

	out_arr_count := entities.Command_count{Values: out_arr, Count: 1}

	return out_arr_count, nil
}
