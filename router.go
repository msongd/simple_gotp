package main

import (
	"embed"
	"io/fs"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

// content is our static web server content.
//go:embed frontend
var staticContent embed.FS

type myFS struct {
	content embed.FS
}

func (c myFS) Open(name string) (fs.File, error) {
	return c.content.Open(path.Join("frontend", name))
}

func MakeRouter(env *Env) *mux.Router {
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
	r.HandleFunc("/static/js/config.js", env.ConfigHandler).Methods("GET")
	if GLOBAL_CFG.UseEmbeddedFrontend {
		r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(myFS{staticContent}))))
	} else {
		r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(GLOBAL_CFG.FrontendDir))))
	}

	r.PathPrefix("/").HandlerFunc(env.CatchAllHandler)
	return r
}
