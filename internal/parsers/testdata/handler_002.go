package handler

import "net/http"

func GET(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func POST(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func PATCH(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func DELETE(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func PUT(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
