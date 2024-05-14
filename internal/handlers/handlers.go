package handlers

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"goshell/internal/services"

	"github.com/gorilla/mux"
)

// func HandlerPostExec(w http.ResponseWriter, r *http.Request) - загрузка скрипта и его выполнение
func HandlerPostExec(w http.ResponseWriter, r *http.Request) {
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

	err = services.CommExec(id)
	if err != nil {
		log.Println(err.Error(), "CommExec error")
	}

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

// func HandlerExecOne(w http.ResponseWriter, r *http.Request) - выполнение скрипта по id
func HandlerExecOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	services.CommExec(i)

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

// func HandlerExec(w http.ResponseWriter, r *http.Request) - выполнение списка скриптов
func HandlerExec(w http.ResponseWriter, r *http.Request) {
	starr := r.URL.Query().Get("ids")
	arr := strings.Split(starr, ",")

	ids := []int{}

	for _, s := range arr {
		i, err := strconv.Atoi(s)
		if err == nil {
			ids = append(ids, i)
		}
	}

	for _, id := range ids {
		go services.CommExec(id)
	}

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

// func HandlerList(w http.ResponseWriter, r *http.Request) - вывод списка команд
func HandlerList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	out_arr_count, err := services.CommGetList(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_count, err := json.Marshal(out_arr_count)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// func HandlerGetOne(w http.ResponseWriter, r *http.Request) - вывод команды по id
func HandlerGetOne(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr_count, err := services.CommGetOne(ctx, i)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_count, err := json.Marshal(out_arr_count)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}

// func HandlerResults(w http.ResponseWriter, r *http.Request) - вывод списка результатов
func HandlerResults(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	out_arr_count, err := services.ResultGetList(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_count, err := json.Marshal(out_arr_count)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// func HealthCheckHandler(w http.ResponseWriter, r *http.Request)
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)

	return
}
