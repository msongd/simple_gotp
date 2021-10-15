package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Listen string `json:"listen"`
	LogDir string `json:"log_dir"`
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
