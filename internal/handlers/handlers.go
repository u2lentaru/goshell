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

	"goshell/internal/services"

	"github.com/gorilla/mux"
)

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

	err = services.CommExec(id)
	if err != nil {
		log.Println(err.Error(), "CommExec error")
	}

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

// func HandleGetOne(w http.ResponseWriter, r *http.Request) - вывод команды по id
func HandleGetOne(w http.ResponseWriter, r *http.Request) {
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

// func HandleResults(w http.ResponseWriter, r *http.Request) - вывод списка результатов
func HandleResults(w http.ResponseWriter, r *http.Request) {
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
