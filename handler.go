package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type AuthContext struct {
	Username string
	IsAdmin  bool
}

func (env *Env) IsUserAllowAccess(username string, r *http.Request) bool {
	authUsername := r.Context().Value("Username").(string)
	isAdmin := r.Context().Value("isAdmin").(bool)

	if !env.Cfg.NoAuth && isAdmin && authUsername != "" && authUsername != username {
		// must auth, not admin, already logged, & access different user url
		return false
	}
	return true
}

func (env *Env) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	o := &OTPVerifyRequest{}
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		log.Println(err)
		http.Error(w, "Error decoding response object", http.StatusBadRequest)
		return
	}
	u, found := env.Db.GetActiveTokenURL(o.Username)
	if !found {
		log.Println("Not found active token for user:", o.Username)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	k, err := otp.NewKeyFromURL(u)
	if err != nil {
		log.Println("Invalid token url:", u)
		http.Error(w, "Invalid token url", http.StatusInternalServerError)
		return
	}
	valid := totp.Validate(o.OTP, k.Secret())
	fmt.Fprintf(w, "{\"valid\":%t}", valid)
}

func (env *Env) CatchAllHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "In Home\n")
	log.Println("CatchAll:", r.RequestURI)
	http.Redirect(w, r, "/static/index.html", http.StatusTemporaryRedirect)
}

func (env *Env) GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("In GetAllUserHandler()")
	var allUsers []*UserDetail
	username := r.Context().Value("Username").(string)
	isAdmin := r.Context().Value("isAdmin").(bool)
	log.Printf("\n[ctx] u:%s a:%t", username, isAdmin)
	if !env.Cfg.NoAuth && !isAdmin && username != "" {
		user, ok := env.Db.Get(username)
		if !ok {
			log.Println("Get user detail err for user:", username)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		allUsers = make([]*UserDetail, 1)
		allUsers[0] = user
	} else {
		log.Println("Admin or noauth access, get all users")
		allUsers = env.Db.GetAllUsers()
	}
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
	username := r.Context().Value("Username").(string)
	isAdmin := r.Context().Value("isAdmin").(bool)

	if !env.Cfg.NoAuth && !isAdmin {
		log.Println("Regular user is not allow to add user:", username)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
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
	if !env.IsUserAllowAccess(username, r) {
		log.Println("Logged user is not allow to get detail for user:", username)
		http.Error(w, "Not found", http.StatusNotFound)
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
	isAdmin := r.Context().Value("isAdmin").(bool)
	authUsername := r.Context().Value("Username").(string)

	username, found := vars["user"]
	if !found {
		log.Println("Not found username in uri")
		http.Error(w, "Error getting username in uri", http.StatusBadRequest)
		return
	}
	if !env.Cfg.NoAuth && !isAdmin {
		log.Println("Regular user is not allow to delete user:", authUsername)
		http.Error(w, "Not found", http.StatusNotFound)
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
	if !env.IsUserAllowAccess(username, r) {
		log.Println("Logged user is not allow to update for user:", username)
		http.Error(w, "Not found", http.StatusNotFound)
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
	if !env.IsUserAllowAccess(username, r) {
		log.Println("Logged user is not allow to get tokens for user:", username)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	// mask all url
	for i, _ := range tokens {
		maskedUrl, err := url.Parse(tokens[i].URL)
		if err != nil {
			log.Println("Parsing url err:", err)
		}
		//log.Println(maskedUrl.Scheme)
		//log.Println(maskedUrl.Host)
		//log.Println(maskedUrl.Path)
		tokens[i].URL = fmt.Sprintf("%s://%s%s***", maskedUrl.Scheme, maskedUrl.Host, maskedUrl.Path)
	}
	//
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

func (env *Env) ImportTokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "In Home\n")
	username, found := vars["user"]
	if !found {
		log.Println("Not found username in uri")
		http.Error(w, "Error getting username in uri", http.StatusBadRequest)
		return
	}
	if !env.IsUserAllowAccess(username, r) {
		log.Println("Logged user is not allow to import token for user:", username)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	u := &struct {
		Url string
	}{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(err)
		http.Error(w, "Error decoding response object", http.StatusBadRequest)
		return
	}
	env.Db.ImportToken(username, u.Url)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "{}")
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
	if !env.IsUserAllowAccess(username, r) {
		log.Println("Logged user is not allow to add token for user:", username)
		http.Error(w, "Not found", http.StatusNotFound)
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
	w.WriteHeader(http.StatusCreated)
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
	if !env.IsUserAllowAccess(username, r) {
		log.Println("Logged user is not allow to get token detail for user:", username)
		http.Error(w, "Not found", http.StatusNotFound)
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
	vars := mux.Vars(r)
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
	if !env.IsUserAllowAccess(username, r) {
		log.Println("Logged user is not allow to get token qr for user:", username)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	token, err := env.Db.GetToken(username, tokenId)
	if err != nil {
		log.Println("Get token err:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	key, err := otp.NewKeyFromURL(token.URL)
	if err != nil {
		log.Println("Create token obj err:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	qrImg, err := key.Image(200, 200)
	if err != nil {
		log.Println("Create qr img err:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	png.Encode(&buf, qrImg)
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	dataUrl := fmt.Sprintf("data:image/png;base64,%s", encoded)
	returnedImage := &struct{ Img string }{Img: dataUrl}
	b, err := json.Marshal(returnedImage)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
func (env *Env) GetOTPHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "In Home\n")
}

func (env *Env) GetAllOTPHandler(w http.ResponseWriter, r *http.Request) {
	var allUserOTPs []*UserDetail
	username := r.Context().Value("Username").(string)
	isAdmin := r.Context().Value("isAdmin").(bool)

	if !env.Cfg.NoAuth && !isAdmin && username != "" {
		user, ok := env.Db.Get(username)
		if !ok {
			log.Println("No user:", username)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		allUserOTPs = make([]*UserDetail, 1)
		allUserOTPs[0] = user
	} else {
		log.Println("Admin or noauth access, get all users")
		allUserOTPs = env.Db.GetAllUserOTPs()
	}

	b, err := json.Marshal(allUserOTPs)
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
	if !env.IsUserAllowAccess(username, r) {
		log.Println("Logged user is not allow to update token for user:", username)
		http.Error(w, "Not found", http.StatusNotFound)
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
	if !env.IsUserAllowAccess(username, r) {
		log.Println("Logged user is not allow to get token qr for user:", username)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	ok, _ := env.Db.RemoveToken(username, tokenId)
	if !ok {
		log.Println("Not remove token", tokenId, "for", username)
		http.Error(w, "Token not found", http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{}")
}

func (env *Env) ConfigHandler(w http.ResponseWriter, r *http.Request) {
	var cfg string
	if env.Cfg.KeycloakCfg == nil || env.Cfg.KeycloakCfg.AuthUrl == "" {
		cfg = ""
	} else {
		cfgTemplate2 := `var keycloakConfig = {
			url: '%s',
			realm: '%s',
			clientId: '%s'
		};
		var ADMIN_ROLE = "%s";`
		cfg = fmt.Sprintf(cfgTemplate2, env.Cfg.KeycloakCfg.AuthUrl, env.Cfg.KeycloakCfg.Realm, env.Cfg.KeycloakCfg.ClientId, env.Cfg.KeycloakCfg.ClaimRoleAdmin)
	}
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", cfg)
}

func (env *Env) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if env.Cfg.NoAuth {
			next.ServeHTTP(w, r)
			return
		}
		tokenBearer := r.Header.Get("Authorization")
		if tokenBearer == "" {
			log.Println("Authorization header not found for url:", r.URL.String())
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		tokens := strings.Split(tokenBearer, " ")
		if len(tokens) != 2 && tokens[0] != "Bearer" {
			log.Println("Bearer text not found: got:", tokenBearer, "for url:", r.URL.String())
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		token := tokens[1]
		tok, err := VerifyJWTToken(token, env.Cfg.KeycloakCfg.JwkUrl, env.Cfg.KeycloakCfg.JwkUrlInsecure)
		if err != nil {
			log.Println("Verify token err:", err, "for url:", r.URL.String())
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		ok := VerifyClaimISS(tok.Claims, env.Cfg)
		if !ok {
			log.Println("Verify token issuer err for url:", r.URL.String())
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		ok, isAdmin := VerifyRealmRole(env.Cfg.KeycloakCfg.ClaimRealmRole, tok.Claims, env)
		if !ok {
			log.Println("Verify token role err for url:", r.URL.String())
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		username := GetUsernameFromJwt(tok.Claims)
		ctx := context.WithValue(r.Context(), "Username", username)
		ctx = context.WithValue(ctx, "isAdmin", isAdmin)
		//env.Username = GetUsernameFromJwt(tok.Claims)
		//env.IsAdmin = IsAdministrator(tok.Claims, env.Cfg)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
