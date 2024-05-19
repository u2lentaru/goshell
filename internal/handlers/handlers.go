package handlers

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"goshell/internal/entities"
	"goshell/internal/services"

	"github.com/gorilla/mux"
)

// type ifCommandService interface
type ifCommandService interface {
	PostExec(ctx context.Context, bs []byte) error
	ExecOne(ctx context.Context, id int) error
	Exec(ctx context.Context, ids []int) error
	GetList(ctx context.Context) (entities.Command_count, error)
	GetOne(ctx context.Context, i int) (entities.Command_count, error)
}

// type ifResultService interface
type ifResultService interface {
	GetList(ctx context.Context) (entities.Result_count, error)
}

// func HandlerPostExec(w http.ResponseWriter, r *http.Request) - загрузка скрипта и его выполнение
func HandlerPostExec(w http.ResponseWriter, r *http.Request) {
	var esv ifCommandService
	esv = services.NewCommandService()
	ctx := context.Background()
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = esv.PostExec(ctx, body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

// func HandlerExecOne(w http.ResponseWriter, r *http.Request) - выполнение скрипта по id
func HandlerExecOne(w http.ResponseWriter, r *http.Request) {
	var esv ifCommandService
	esv = services.NewCommandService()
	ctx := context.Background()
	vars := mux.Vars(r)

	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	esv.ExecOne(ctx, i)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

// func HandlerExec(w http.ResponseWriter, r *http.Request) - выполнение списка скриптов
func HandlerExec(w http.ResponseWriter, r *http.Request) {
	var esv ifCommandService
	esv = services.NewCommandService()
	ctx := context.Background()
	starr := r.URL.Query().Get("ids")
	arr := strings.Split(starr, ",")

	ids := []int{}

	for _, s := range arr {
		i, err := strconv.Atoi(s)
		if err == nil {
			ids = append(ids, i)
		}
	}

	err := esv.Exec(ctx, ids)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/commands", http.StatusSeeOther)
}

// func HandlerList(w http.ResponseWriter, r *http.Request) - вывод списка команд
func HandlerList(w http.ResponseWriter, r *http.Request) {
	var esv ifCommandService
	esv = services.NewCommandService()
	ctx := context.Background()

	out_arr_count, err := esv.GetList(ctx)
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
	var esv ifCommandService
	esv = services.NewCommandService()
	ctx := context.Background()

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr_count, err := esv.GetOne(ctx, i)
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
	var esv ifResultService
	esv = services.NewResultService()
	ctx := context.Background()

	out_arr_count, err := esv.GetList(ctx)
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
