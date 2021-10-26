package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Listen              string          `json:"listen"`
	SSLKeyFile          string          `json:"ssl_key"`
	SSLCertFile         string          `json:"ssl_cert"`
	LogDir              string          `json:"log_dir"`
	DataFile            string          `json:"data_file"`
	FrontendDir         string          `json:"frontend_dir"`
	UseEmbeddedFrontend bool            `json:"use_embedded_frontend"`
	KeycloakCfg         *KeycloakConfig `json:"keycloak_cfg"`
	NoAuth              bool            `json:"-"`
	TLSConfig           *tls.Config     `json:"-"`
}

type KeycloakConfig struct {
	AuthUrl        string `json:"auth_url"`
	Realm          string `json:"realm"`
	ClientId       string `json:"client_id"`
	Secret         string `json:"secret"`
	JwkUrl         string `json:"jwk_url"`
	ClaimIss       string `json:"claim_iss"`
	ClaimRealmRole string `json:"claim_realm_role"`
	ClaimRoleAdmin string `json:"claim_role_admin"`
}

func NewConfig() *Config {
	m := Config{}
	return &m
}

func (cfg *Config) LoadConfig(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		log.Println("Openfile err:", err)
		return err
	}
	defer f.Close()
	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println("Read file err:", err)
		return err
	}
	err = json.Unmarshal(byteValue, cfg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
