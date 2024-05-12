package handlers

import "net/http"

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Root!"))
	return
}

func HandleExec(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Exec!"))
	return
}

func HandleList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello List!"))
	return
}

func HandleGetOne(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello GetOne!"))
	return
}

func HandleFormTypes(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
	return
}
