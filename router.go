package main

import (
	"github.com/gorilla/mux"
)

func MakeRouter() *mux.Router {
	r := mux.NewRouter()
	// Add your routes as needed
	r.HandleFunc("/public/verify", VerifyHandler).Methods("POST")
	r.HandleFunc("/auth/token", TokenHandler)
	r.HandleFunc("/", HomeHandler)
	return r
}
