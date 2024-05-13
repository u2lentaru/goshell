package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"goshell/internal/pgclient"
	"goshell/internal/structs"

	"github.com/gorilla/mux"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Root!"))
	return
}

// func HandlePostExec(w http.ResponseWriter, r *http.Request) - загрузка скрипта и его выполнение
func HandlePostExec(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	dbpool := pgclient.WDB

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte(" body: "))
	w.Write(body)

	_ = os.MkdirAll("/test", 0777)

	lsout, err := exec.Command("ls", "/").Output()
	w.Write([]byte(" ls /: "))
	w.Write(lsout)

	f, err := os.CreateTemp("/test", "*.sh")
	if err != nil {
		log.Println(err.Error())
	}
	defer f.Close()

	fn := f.Name()
	w.Write([]byte(" f.Name "))
	w.Write([]byte(fn))

	if err := os.WriteFile(fn, body, 0777); err != nil {
		log.Println(err.Error())
	}

	lsout, err = exec.Command("ls", "-l", fn).Output()
	w.Write([]byte(" ls /test/file.sh: "))
	w.Write(lsout)

	stout := ""
	sterr := ""
	cmd := fn

	out, err := exec.Command("bash", cmd).Output()

	if err != nil {
		sterr = fmt.Sprintf("Failed to execute command: %s error %s", cmd, err.Error())
	}

	stout = string(out)

	cid := 0
	err = dbpool.QueryRow(ctx, "insert into commands (id, command_text, script_text) values (default, $1, $2) returning id;", cmd, string(body)).Scan(&cid)

	if err != nil {
		// w.Write([]byte("Failed execute command add!"))
		w.Write([]byte(" dberr: "))
		w.Write([]byte(err.Error()))
		w.Write([]byte(" stout: "))
		w.Write([]byte(stout))
		w.Write([]byte(" sterr: "))
		w.Write([]byte(sterr))
		return
	}

	rid := 0
	t := time.Now()
	err = dbpool.QueryRow(ctx, "insert into results (id, id_command, output, time) values (default, $1, $2, $3) returning id;", cid, stout, t.Format("2006-01-02 15:04:05")).Scan(&rid)

	if err != nil {
		// w.Write([]byte("Failed execute command add!"))
		w.Write([]byte(" dberr: "))
		w.Write([]byte(err.Error()))
		w.Write([]byte(" stout: "))
		w.Write([]byte(stout))
		w.Write([]byte(" sterr: "))
		w.Write([]byte(sterr))
		return
	}

	w.Write([]byte(" stout: "))
	w.Write([]byte(stout))
	w.Write([]byte(" sterr: "))
	w.Write([]byte(sterr))

	return
}

// func HandleExecOne(w http.ResponseWriter, r *http.Request) - выполнение скрипта по id
func HandleExecOne(w http.ResponseWriter, r *http.Request) {
	// db
	ctx := context.Background()

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	dbpool := pgclient.WDB
	g := structs.Command{}

	err = dbpool.QueryRow(ctx, "SELECT * from public.commands where id=$1;", i).Scan(&g.Id, &g.CommandText, &g.ScriptText)

	if err != nil {
		log.Println(err.Error(), "exec_one")
		http.Redirect(w, r, "/results", http.StatusSeeOther)
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
		w.Write([]byte("Failed execute command add!"))
		return
	}

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

// func HandleExec(w http.ResponseWriter, r *http.Request) - выполнение списка скриптов
func HandleExec(w http.ResponseWriter, r *http.Request) {
	// db
	return
}

// func HandleList(w http.ResponseWriter, r *http.Request) - вывод списка команд
func HandleList(w http.ResponseWriter, r *http.Request) {
	// db
	ctx := context.Background()

	dbpool := pgclient.WDB
	gs := structs.Command{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT count(*) from public.commands;").Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "commands_count")
		return
	}

	out_arr := make([]structs.Command, 0, gsc)

	rows, err := dbpool.Query(ctx, "SELECT * from public.commands;")
	if err != nil {
		log.Println(err.Error(), "commands_list")
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.CommandText, &gs.ScriptText)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_arr_count := structs.Command_count{Values: out_arr, Count: gsc}

	// handler
	out_count, err := json.Marshal(out_arr_count)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// func HandleGetOne(w http.ResponseWriter, r *http.Request) - вывод команды по id
func HandleGetOne(w http.ResponseWriter, r *http.Request) {
	// db
	ctx := context.Background()

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	dbpool := pgclient.WDB
	out_arr := []structs.Command{}
	g := structs.Command{}

	err = dbpool.QueryRow(ctx, "SELECT * from public.commands where id=$1;", i).Scan(&g.Id, &g.CommandText, &g.ScriptText)

	if err != nil {
		log.Println(err.Error(), "commands_one")
		out_arr_count := structs.Command_count{Values: []structs.Command{}, Count: 0}

		out_count, err := json.Marshal(out_arr_count)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write(out_count)

		return
	}

	out_arr = append(out_arr, g)

	out_arr_count := structs.Command_count{Values: out_arr, Count: 1}

	// handler
	out_count, err := json.Marshal(out_arr_count)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}

// func HandleResults(w http.ResponseWriter, r *http.Request) - вывод списка результатов
func HandleResults(w http.ResponseWriter, r *http.Request) {
	// db
	ctx := context.Background()

	dbpool := pgclient.WDB
	gs := structs.Result{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT count(*) from public.results;").Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "results_count")
		return
	}

	out_arr := make([]structs.Result, 0, gsc)

	rows, err := dbpool.Query(ctx, "SELECT id, id_command, output, time::text as ts from public.results;")
	if err != nil {
		log.Println(err.Error(), "results_list")
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.IdC, &gs.Output, &gs.TS)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	out_arr_count := structs.Result_count{Values: out_arr, Count: gsc}

	// handler
	out_count, err := json.Marshal(out_arr_count)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}
