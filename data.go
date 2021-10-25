package main

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TokenDetail struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Issuer string `json:"issuer"`
}

type UserDetail struct {
	Username    string        `json:"username"`
	ActiveToken string        `json:"active_token"`
	Tokens      []TokenDetail `json:"tokens"`
	Total       int           `json:"total"`
	CurrentCode string        `json:"current_code"`
}

type OTPVerifyRequest struct {
	Username string `json:"username"`
	OTP      string `json:"otp"`
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

func (u *UserDetail) Cloned(nowTime time.Time) *UserDetail {
	newUser := NewUser()
	newUser.ActiveToken = u.ActiveToken
	newUser.Username = u.Username
	newUser.Tokens = make([]TokenDetail, len(u.Tokens))
	//copy(newUser.Tokens, u.Tokens)
	for i, t := range u.Tokens {
		newUser.Tokens[i].ID = t.ID
		k, err := otp.NewKeyFromURL(t.URL)
		if err != nil {
			log.Println("Parsing url to token err:", err)
		} else {
			newUser.Tokens[i].Issuer = k.Issuer()
		}
		if t.ID == u.ActiveToken {
			log.Println("Generate code for:user:", u.Username, ":time:", nowTime)
			code, err := totp.GenerateCode(k.Secret(), nowTime)
			if err != nil {
				log.Println("Error getting code:", err)
			}
			newUser.CurrentCode = code
		}
	}
	newUser.Total = len(u.Tokens)
	return newUser
}

func (cfg *OtpConfig) Get(username string) (*UserDetail, bool) {
	cfg.Lock()
	defer cfg.Unlock()
	c, found := cfg.Users[username]
	if !found {
		return nil, false
	}
	nowTime := time.Now()
	cloned := c.Cloned(nowTime)
	return cloned, found
}

func (cfg *OtpConfig) GetActiveTokenURL(username string) (string, bool) {
	cfg.Lock()
	defer cfg.Unlock()
	c, found := cfg.Users[username]
	if !found {
		return "", false
	}
	for _, t := range c.Tokens {
		if c.ActiveToken == t.ID {
			return t.URL, true
		}
	}
	return "", false
}

func (cfg *OtpConfig) Remove(username string) bool {
	cfg.Lock()
	defer cfg.Unlock()
	_, found := cfg.Users[username]
	if found {
		delete(cfg.Users, username)
		return true
	}
	return false
}

func (cfg *OtpConfig) GetAllUsers() []*UserDetail {
	cfg.Lock()
	defer cfg.Unlock()
	all := make([]*UserDetail, len(cfg.Users))
	i := 0
	nowTime := time.Now()
	for _, v := range cfg.Users {
		all[i] = v.Cloned(nowTime)
		i = i + 1
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].Username < all[j].Username
	})

	return all
}
func (cfg *OtpConfig) GetAllUserOTPs() []*UserDetail {
	return cfg.GetAllUsers()
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
	return found
}

func (cfg *OtpConfig) GetAllTokens(username string) ([]TokenDetail, error) {
	cfg.Lock()
	defer cfg.Unlock()
	u, found := cfg.Users[username]
	if !found {
		return nil, nil
	}
	tokens := make([]TokenDetail, len(u.Tokens))
	for i, t := range u.Tokens {
		tokens[i] = t
		k, err := otp.NewKeyFromURL(t.URL)
		if err != nil {
			log.Println("Parsing url to token err:", err)
		} else {
			tokens[i].Issuer = k.Issuer()
		}
	}
	sort.Slice(tokens, func(i, j int) bool { return tokens[i].ID < tokens[j].ID })
	//copy(tokens, u.Tokens)
	return tokens, nil
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

	newKey, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: username,
	})

	if err != nil {
		log.Println("Generate new token error:", err)
		return false, err
	}
	//t.ID = fmt.Sprintf("%d", l+1)
	t := TokenDetail{}
	t.URL = newKey.URL()
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
func (cfg *OtpConfig) ImportToken(username string, url string) (bool, error) {
	cfg.Lock()
	defer cfg.Unlock()
	u, found := cfg.Users[username]
	if !found {
		u = NewUser()
		cfg.Users[username] = u
		u.Username = username
	}

	newKey, err := otp.NewKeyFromURL(url)
	if err != nil {
		log.Println("Generate new token error:", err)
		return false, err
	}
	//t.ID = fmt.Sprintf("%d", l+1)
	t := TokenDetail{}
	t.URL = newKey.URL()
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

func (cfg *OtpConfig) GetToken(username string, tokenId string) (*TokenDetail, error) {
	cfg.Lock()
	defer cfg.Unlock()
	u, found := cfg.Users[username]
	if !found {
		return nil, nil
	}
	for _, t := range u.Tokens {
		if tokenId == t.ID {
			returned := t
			k, err := otp.NewKeyFromURL(t.URL)
			if err != nil {
				log.Println("Parsing url to token err:", err)
				return nil, err
			}
			returned.Issuer = k.Issuer()
			return &returned, nil
		}
	}
	return nil, nil
}

func (cfg *OtpConfig) UpdateToken(username string, tokenId string, url string) (bool, error) {
	cfg.Lock()
	defer cfg.Unlock()
	u, found := cfg.Users[username]
	if !found {
		u = NewUser()
		cfg.Users[username] = u
		u.Username = username
	}
	key, err := otp.NewKeyFromURL(url)
	if err != nil {
		log.Println("Parse url error:", err)
		return false, err
	}
	keyAccName := key.AccountName()
	if keyAccName != username {
		log.Println("Token url has different username:", keyAccName)
		return false, nil
	}
	for i, t := range u.Tokens {
		if tokenId == t.ID {
			u.Tokens[i].URL = url
			return true, nil
		}
	}
	return false, nil
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

func (tok TokenDetail) GetOtp() {

}

func Rand20() []byte {
	c := 20
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
