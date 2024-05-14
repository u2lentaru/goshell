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
func CommExec(id int) {
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
		return
	}

	return
}

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
