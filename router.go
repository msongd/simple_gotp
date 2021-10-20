package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func MakeRouter(env *Env, staticDir string) *mux.Router {
	r := mux.NewRouter()
	// Add your routes as needed
	r.HandleFunc("/public/verify", env.VerifyHandler).Methods("POST")
	r.HandleFunc("/auth/user/{user}", env.GetUserHandler).Methods("GET")
	r.HandleFunc("/auth/user/{user}", env.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/auth/user/{user}", env.UpdateUserHandler).Methods("POST")
	r.HandleFunc("/auth/user", env.GetAllUserHandler).Methods("GET")
	r.HandleFunc("/auth/user", env.AddUserHandler).Methods("POST")
	r.HandleFunc("/auth/qr/{user}/{token}", env.GetTokenQRHandler).Methods("POST")
	r.HandleFunc("/auth/otp/{user}/{token}", env.GetOTPHandler).Methods("GET")
	r.HandleFunc("/auth/otp", env.GetAllOTPHandler).Methods("GET")
	r.HandleFunc("/auth/token/{user}/import", env.ImportTokenHandler).Methods("POST")
	//r.HandleFunc("/auth/token/{user}/{token}", env.GetTokenHandler).Methods("GET")
	//r.HandleFunc("/auth/token/{user}/{token}", env.UpdateTokenHandler).Methods("POST")
	r.HandleFunc("/auth/token/{user}/{token}", env.DeleteTokenHandler).Methods("DELETE")
	r.HandleFunc("/auth/token/{user}", env.GetAllTokenHandler).Methods("GET")
	r.HandleFunc("/auth/token/{user}", env.AddTokenHandler).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	r.PathPrefix("/").HandlerFunc(env.CatchAllHandler)
	return r
}
