package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"goshell/internal/entities"
	"goshell/internal/pgclient"
	"goshell/internal/services"

	"github.com/gorilla/mux"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Root!"))
	return
}

// func HandlePostExec(w http.ResponseWriter, r *http.Request) - загрузка скрипта и его выполнение
func HandlePostExec(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id, err := services.CommSave(body)
	if err != nil {
		log.Println(err.Error(), "CommSave error")
		return
	}

	services.CommExec(id)

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

// func HandleExecOne(w http.ResponseWriter, r *http.Request) - выполнение скрипта по id
func HandleExecOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	services.CommExec(i)

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

// func HandleExec(w http.ResponseWriter, r *http.Request) - выполнение списка скриптов
func HandleExec(w http.ResponseWriter, r *http.Request) {
	a := []int{1, 2, 3}

	for _, i := range a {
		go services.CommExec(i)
	}

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

// func HandleList(w http.ResponseWriter, r *http.Request) - вывод списка команд
func HandleList(w http.ResponseWriter, r *http.Request) {
	// db
	ctx := context.Background()

	dbpool := pgclient.WDB
	gs := entities.Command{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT count(*) from public.commands;").Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "commands_count")
		return
	}

	out_arr := make([]entities.Command, 0, gsc)

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

	out_arr_count := entities.Command_count{Values: out_arr, Count: gsc}

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
	out_arr := []entities.Command{}
	g := entities.Command{}

	err = dbpool.QueryRow(ctx, "SELECT * from public.commands where id=$1;", i).Scan(&g.Id, &g.CommandText, &g.ScriptText)

	if err != nil {
		log.Println(err.Error(), "commands_one")
		out_arr_count := entities.Command_count{Values: []entities.Command{}, Count: 0}

		out_count, err := json.Marshal(out_arr_count)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write(out_count)

		return
	}

	out_arr = append(out_arr, g)

	out_arr_count := entities.Command_count{Values: out_arr, Count: 1}

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
	gs := entities.Result{}

	gsc := 0
	err := dbpool.QueryRow(ctx, "SELECT count(*) from public.results;").Scan(&gsc)

	if err != nil {
		log.Println(err.Error(), "results_count")
		return
	}

	out_arr := make([]entities.Result, 0, gsc)

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

	out_arr_count := entities.Result_count{Values: out_arr, Count: gsc}

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
func HandleTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	arr := r.URL.Query().Get("id")

	var url, err = url.ParseQuery("a=1&a=2&b=3")
	// var url, err = url.ParseQuery(r.URL.Query())

	if err != nil {
		log.Println("failed to parse:", err)
	}

	fmt.Println(url["a"])

	w.Write([]byte(" vars[id] "))
	w.Write([]byte(vars["id"]))

	w.Write([]byte(" arr: "))
	w.Write([]byte(arr))

	// url.Get("id")
	return
}
