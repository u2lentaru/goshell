package handlers

import "net/http"

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
	return
}

func HandleFormTypes(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
	return
}
