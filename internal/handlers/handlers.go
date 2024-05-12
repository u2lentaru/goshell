package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"goshell/internal/pgclient"
	"goshell/internal/structs"

	"github.com/gorilla/mux"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Root!"))
	return
}

func HandleExec(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Exec!"))
	return
}

// func HandleList(w http.ResponseWriter, r *http.Request)
func HandleList(w http.ResponseWriter, r *http.Request) {
	// db
	ctx := context.Background()

	dbpool := pgclient.WDB
	gs := structs.Command{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT count(*) from commands;").Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "commands_count")
		return
	}

	out_arr := make([]structs.Command, 0, gsc)

	rows, err := dbpool.Query(ctx, "SELECT * from commands;")
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

	err = dbpool.QueryRow(ctx, "SELECT * from commands where id=$1;", i).Scan(&g.Id, &g.CommandText, &g.ResultText)

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
