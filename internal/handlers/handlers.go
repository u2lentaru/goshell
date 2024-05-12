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

	"goshell/internal/pgclient"
	"goshell/internal/structs"

	"github.com/gorilla/mux"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Root!"))
	return
}

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

	if err := os.WriteFile("/test/file.sh", body, 0777); err != nil {
		log.Println(err.Error())
	}

	lsout, err = exec.Command("ls", "-l", "/test/file.sh").Output()
	w.Write([]byte(" ls /test/file.sh: "))
	w.Write(lsout)

	stout := ""
	sterr := ""
	cmd := "/test/file.sh"

	out, err := exec.Command("bash", cmd).Output()

	if err != nil {
		sterr = fmt.Sprintf("Failed to execute command: %s error %s", cmd, err.Error())
	}

	stout = string(out)

	dbres, err := dbpool.Exec(ctx, "insert into commands (id, command_text, result_text) values (default, $1, $2);", cmd, stout)

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

	w.Write([]byte(" dbres: "))
	w.Write([]byte(dbres))
	w.Write([]byte(" stout: "))
	w.Write([]byte(stout))
	w.Write([]byte(" sterr: "))
	w.Write([]byte(sterr))

	return
}

func HandleExec(w http.ResponseWriter, r *http.Request) {
	// db
	ctx := context.Background()
	dbpool := pgclient.WDB

	stout := ""
	sterr := ""
	// cmd := "cat /proc/cpuinfo | egrep '^model name' | uniq | awk '{print substr($0, index($0,$4))}'"
	cmd := "ls"
	// out, err := exec.Command("bash", "-c", cmd).Output()
	out, err := exec.Command(cmd).Output()
	if err != nil {
		sterr = fmt.Sprintf("Failed to execute command: %s error %s", cmd, err.Error())
	}

	stout = string(out)

	// ai := 0
	// err = dbpool.QueryRow(ctx, "insert into commands (id, command_text, result_text) values (default, $1, $2) returning id into i;", cmd, stout).Scan(&ai)
	dbres, err := dbpool.Exec(ctx, "insert into commands (id, command_text, result_text) values (default, $1, $2);", cmd, stout)

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

	w.Write([]byte(" dbres: "))
	w.Write([]byte(dbres))
	w.Write([]byte(" stout: "))
	w.Write([]byte(stout))
	w.Write([]byte(" sterr: "))
	w.Write([]byte(sterr))

	return
}

// func HandleList(w http.ResponseWriter, r *http.Request)
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
		err = rows.Scan(&gs.Id, &gs.CommandText, &gs.ResultText)
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

// func HandleGetOne(w http.ResponseWriter, r *http.Request)
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

	err = dbpool.QueryRow(ctx, "SELECT * from public.commands where id=$1;", i).Scan(&g.Id, &g.CommandText, &g.ResultText)

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
