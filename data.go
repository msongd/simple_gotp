package main

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
)

type TokenDetail struct {
	ID          string `json:"id"`
	AccountName string `json:"account_name"`
	Issuer      string `json:"issuer"`
	//	Period int `json:"period"`
	Secret string `json:"secret"`
	//	Digits int `json:"digits"`
}

type UserDetail struct {
	Username    string        `json:"username"`
	ActiveToken string        `json:"active_token"`
	Tokens      []TokenDetail `json:"tokens"`
}

type OtpConfig struct {
	sync.RWMutex
	Users map[string]*UserDetail
}

func NewUser() *UserDetail {
	u := UserDetail{}
	u.Tokens = make([]TokenDetail, 0)
	return &u
}

func NewOtpConfig() *OtpConfig {
	m := OtpConfig{}
	m.Users = make(map[string]*UserDetail)
	return &m
}

func (u *UserDetail) Cloned() *UserDetail {
	newUser := NewUser()
	newUser.ActiveToken = u.ActiveToken
	newUser.Username = u.Username
	copy(newUser.Tokens, u.Tokens)
	return newUser
}

func (cfg *OtpConfig) Get(username string) (*UserDetail, bool) {
	cfg.Lock()
	defer cfg.Unlock()
	c, found := cfg.Users[username]
	cloned := c.Cloned()
	return cloned, found
}

func (cfg *OtpConfig) GetAllUsers() []*UserDetail {
	cfg.Lock()
	defer cfg.Unlock()
	all := make([]*UserDetail, len(cfg.Users))
	i := 0
	for _, v := range cfg.Users {
		all[i] = v.Cloned()
		i = i + 1
	}
	return all
}

func (cfg *OtpConfig) AddUser(username string) (bool, error) {
	cfg.Lock()
	defer cfg.Unlock()
	_, found := cfg.Users[username]
	if !found {
		u := NewUser()
		u.Username = username
		cfg.Users[username] = u
		return true, nil
	}
	return false, nil
}

func (cfg *OtpConfig) HasUser(username string) bool {
	cfg.Lock()
	defer cfg.Unlock()
	_, found := cfg.Users[username]
	if !found {
		return false
	}
	return true
}

func (cfg *OtpConfig) AddToken(username string, issuer string) (bool, error) {
	cfg.Lock()
	defer cfg.Unlock()
	u, found := cfg.Users[username]
	if !found {
		u = NewUser()
		cfg.Users[username] = u
		u.Username = username
	}
	t := TokenDetail{}
	t.AccountName = username
	t.Issuer = issuer
	t.Secret = base32.StdEncoding.EncodeToString(Rand16())
	l := len(u.Tokens)
	if l > 0 {
		lastIdStr := u.Tokens[l-1].ID
		id, err := strconv.Atoi(lastIdStr)
		if err != nil {
			log.Println("Convert last id to int err:", err)
			t.ID = fmt.Sprintf("%s1", lastIdStr)
		} else {
			t.ID = fmt.Sprintf("%d", id+1)
		}
	} else {
		t.ID = "1"
	}
	//t.ID = fmt.Sprintf("%d", l+1)
	u.Tokens = append(u.Tokens, t)
	return true, nil
}
func (cfg *OtpConfig) UpdateToken(username string, tokenId string, issuer string, secret string) (bool, error) {
	cfg.Lock()
	defer cfg.Unlock()
	u, found := cfg.Users[username]
	if !found {
		u = NewUser()
		cfg.Users[username] = u
		u.Username = username
	}
	for i, t := range u.Tokens {
		if tokenId == t.ID {
			u.Tokens[i].AccountName = username
			u.Tokens[i].Issuer = issuer
			u.Tokens[i].Secret = secret
			return true, nil
		}
	}
	return true, nil
}

func (cfg *OtpConfig) RemoveToken(username string, tokenId string) (bool, error) {
	cfg.Lock()
	defer cfg.Unlock()
	u, found := cfg.Users[username]
	if !found {
		return false, nil
	}
	for i, t := range u.Tokens {
		if tokenId == t.ID {
			copy(u.Tokens[i:], u.Tokens[i+1:])        // Shift a[i+1:] left one index.
			u.Tokens[len(u.Tokens)-1] = TokenDetail{} // Erase last element (write zero value).
			u.Tokens = u.Tokens[:len(u.Tokens)-1]     // Truncate slice.
			if u.ActiveToken == tokenId {
				u.ActiveToken = ""
			}
			return true, nil
		}
	}
	return false, nil
}

func (cfg *OtpConfig) SetActiveToken(username string, tokenId string) (bool, error) {
	cfg.Lock()
	defer cfg.Unlock()
	u, found := cfg.Users[username]
	if !found {
		return false, nil
	}
	for _, t := range u.Tokens {
		if tokenId == t.ID {
			u.ActiveToken = tokenId
			return true, nil
		}
	}
	return false, nil
}

func LoadFromFile(filename string) (*OtpConfig, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Println("Openfile err:", err)
		return nil, err
	}
	defer f.Close()
	cfg := NewOtpConfig()
	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println("Read file err:", err)
		return nil, err
	}
	json.Unmarshal(byteValue, cfg)
	return cfg, nil
}
func (cfg *OtpConfig) SaveToFile(filename string) error {
	log.Println("Saving data...")
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Println(err)
		return err
	}
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		log.Println("Save data to file err:", err)
		return err
	}
	log.Println("Data saved...")
	return nil
}

func (cfg *OtpConfig) Dump() {
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(b))
}

func (cfg *OtpConfig) GetActiveToken(username string) (bool, *TokenDetail) {
	u, found := cfg.Users[username]
	if !found {
		return false, nil
	}
	for _, t := range u.Tokens {
		if u.ActiveToken == t.ID {
			cloned := t
			return true, &cloned
		}
	}
	return false, nil
}

func Rand16() []byte {
	c := 16
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("error:", err)
		return nil
	}
	return b
}

func MakeSecret(b []byte) string {
	return base32.StdEncoding.EncodeToString(b)
}
