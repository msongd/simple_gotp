package main

import (
	"embed"
	"io/fs"
	"log"
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
	log.Printf("%+v\n", env.Cfg)
	r := mux.NewRouter()
	//r.Use(env.AuthenticationMiddleware)
	// Add your routes as needed
	r.HandleFunc("/public/verify", env.VerifyHandler).Methods("POST")

	authRoute := r.PathPrefix("/auth").Subrouter()
	authRoute.Use(env.AuthenticationMiddleware)
	authRoute.HandleFunc("/user/{user}/alias/{alias}", env.RemoveAliasHandler).Methods("DELETE")
	authRoute.HandleFunc("/user/{user}/alias", env.AddAliasHandler).Methods("POST")
	authRoute.HandleFunc("/user/{user}", env.GetUserHandler).Methods("GET")
	authRoute.HandleFunc("/user/{user}", env.DeleteUserHandler).Methods("DELETE")
	authRoute.HandleFunc("/user/{user}", env.UpdateUserHandler).Methods("POST")
	authRoute.HandleFunc("/user", env.GetAllUserHandler).Methods("GET")
	authRoute.HandleFunc("/user", env.AddUserHandler).Methods("POST")
	authRoute.HandleFunc("/qr/{user}/{token}", env.GetTokenQRHandler).Methods("POST")
	authRoute.HandleFunc("/otp/{user}/{token}", env.GetOTPHandler).Methods("GET")
	authRoute.HandleFunc("/otp", env.GetAllOTPHandler).Methods("GET")
	authRoute.HandleFunc("/token/{user}/import", env.ImportTokenHandler).Methods("POST")
	authRoute.HandleFunc("/token/{user}/{token}", env.DeleteTokenHandler).Methods("DELETE")
	authRoute.HandleFunc("/token/{user}", env.GetAllTokenHandler).Methods("GET")
	authRoute.HandleFunc("/token/{user}", env.AddTokenHandler).Methods("POST")

	r.HandleFunc("/static/js/config.js", env.ConfigHandler).Methods("GET")
	if env.Cfg.UseEmbeddedFrontend {
		r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(myFS{staticContent}))))
	} else {
		r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(env.Cfg.FrontendDir))))
	}

	r.PathPrefix("/").HandlerFunc(env.CatchAllHandler)
	return r
}
