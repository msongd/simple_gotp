package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (env *Env) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "In Verify\n")
}

func (env *Env) CatchAllHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "In Home\n")

}

func (env *Env) GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
	allUsers := env.Db.GetAllUsers()
	b, err := json.Marshal(allUsers)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "In Home\n")
	w.Write(b)
}

func (env *Env) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "In Home\n")
	u := NewUser()
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err)
		http.Error(w, "Error decoding response object", http.StatusBadRequest)
		return
	}
	ok, err := env.Db.AddUser(u.Username)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding response object", http.StatusBadRequest)
		return
	}
	if ok {
		// added
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "{}")
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "{}")
	}
}

func (env *Env) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, found := vars["user"]
	if !found {
		log.Println("Not found username in uri")
		http.Error(w, "Error getting username in uri", http.StatusBadRequest)
		return
	}
	u, found := env.Db.Get(username)
	if !found {
		log.Println("User not found:", username)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	b, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (env *Env) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "In Home\n")
	username, found := vars["user"]
	if !found {
		log.Println("Not found username in uri")
		http.Error(w, "Error getting username in uri", http.StatusBadRequest)
		return
	}
	found = env.Db.Remove(username)
	if !found {
		log.Println("User not found:", username)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{}")
}
func (env *Env) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "In Home\n")
	username, found := vars["user"]
	if !found {
		log.Println("Not found username in uri")
		http.Error(w, "Error getting username in uri", http.StatusBadRequest)
		return
	}
	u := NewUser()
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err)
		http.Error(w, "Error decoding response object", http.StatusBadRequest)
		return
	}
	if u.ActiveToken == "" {
		log.Println("No ActiveToken info")
		http.Error(w, "Empty active token info", http.StatusBadRequest)
		return
	}
	ok, _ := env.Db.SetActiveToken(username, u.ActiveToken)
	if !ok {
		log.Println("Unable to set active token for ", username)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{}")
}
func (env *Env) GetAllTokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "In Home\n")
	username, found := vars["user"]
	if !found {
		log.Println("Not found username in uri")
		http.Error(w, "Error getting username in uri", http.StatusBadRequest)
		return
	}
	tokens, _ := env.Db.GetAllTokens(username)
	if tokens == nil {
		log.Println("Error getting token for user ", username)
		http.Error(w, "Error getting token", http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(tokens)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
func (env *Env) AddTokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "In Home\n")
	username, found := vars["user"]
	if !found {
		log.Println("Not found username in uri")
		http.Error(w, "Error getting username in uri", http.StatusBadRequest)
		return
	}
	u := &struct{ Issuer string }{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err)
		http.Error(w, "Error decoding response object", http.StatusBadRequest)
		return
	}
	env.Db.AddToken(username, u.Issuer)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{}")
}
func (env *Env) GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "In Home\n")
	username, found := vars["user"]
	if !found {
		log.Println("Not found username in uri")
		http.Error(w, "Error getting username in uri", http.StatusBadRequest)
		return
	}
	tokenId, found := vars["token"]
	if !found {
		log.Println("Not found token id in uri")
		http.Error(w, "Error getting token id in uri", http.StatusBadRequest)
		return
	}
	token, _ := env.Db.GetToken(username, tokenId)
	if token == nil {
		log.Println("Not found token id")
		http.Error(w, "Token not found", http.StatusNotFound)
		return
	}
	b, err := json.Marshal(token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
func (env *Env) GetTokenQRHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "In Home\n")

}
func (env *Env) GetOTPHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "In Home\n")
}

func (env *Env) UpdateTokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "In Home\n")
	username, found := vars["user"]
	if !found {
		log.Println("Not found username in uri")
		http.Error(w, "Error getting username in uri", http.StatusBadRequest)
		return
	}
	tokenId, found := vars["token"]
	if !found {
		log.Println("Not found token id in uri")
		http.Error(w, "Error getting token id in uri", http.StatusBadRequest)
		return
	}
	u := &TokenDetail{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err)
		http.Error(w, "Error decoding response object", http.StatusBadRequest)
		return
	}
	if u.ID != "" && u.ID != tokenId {
		log.Println("Posted token id differ from token id in uri")
		http.Error(w, "Posted token id differ from token id in uri", http.StatusBadRequest)
		return
	}
	env.Db.UpdateToken(username, tokenId, u.URL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{}")
}

func (env *Env) DeleteTokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "In Home\n")
	username, found := vars["user"]
	if !found {
		log.Println("Not found username in uri")
		http.Error(w, "Error getting username in uri", http.StatusBadRequest)
		return
	}
	tokenId, found := vars["token"]
	if !found {
		log.Println("Not found token id in uri")
		http.Error(w, "Error getting token id in uri", http.StatusBadRequest)
		return
	}
	ok, _ := env.Db.RemoveToken(username, tokenId)
	if !ok {
		log.Println("Not remove token", tokenId, "for", username)
		http.Error(w, "Token not found", http.StatusNotFound)
	}
}
