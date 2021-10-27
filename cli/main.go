package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

//NAS_IDENTIFIER="Vendor"
//USER_NAME="tungdo"
//PWD=/
//NAS_IP_ADDRESS=10.1.5.17
//MESSAGE_AUTHENTICATOR=0x7b0c9f30114c84f217695a11fba597d7
//USER_PASSWORD="123456"
//EVENT_TIMESTAMP="Oct 21 2021 10:34:51 +07"

var DEFAULT_REQUEST_FORMAT = `{"username":%s,"otp":%s}`

type Resp struct {
	Valid bool `json:"valid"`
}

var (
	DEFAULT_OTP_URL  = flag.String("u", "http://127.0.0.1:8080/public/verify", "URL to query")
	DEFAULT_LOG_FILE = flag.String("l", "/tmp/simple_gotp.log", "log file")
	IS_DEBUG         = flag.Bool("d", false, "Turn on more debug messages")
)

func main() {
	if len(os.Args) == 1 {
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(&lumberjack.Logger{
		Filename:   *DEFAULT_LOG_FILE,
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	})

	//compatible with radius exec module
	username := os.Getenv("USER_NAME")
	password := os.Getenv("USER_PASSWORD")
	client_ip := os.Getenv("NAS_IP_ADDRESS")
	nas_id := os.Getenv("NAS_IDENTIFIER")
	timestamp := os.Getenv("EVENT_TIMESTAMP")

	log.Printf("[info] [username:%s] [client_ip:%s] [nas_id:%s] [time:%s]\n", username, client_ip, nas_id, timestamp)

	ok, err := verify(DEFAULT_OTP_URL, username, password)
	if err != nil {
		//#  = 2 : fail      the module failed
		os.Exit(1)
	}
	if !ok {
		os.Exit(1)
	}
	os.Exit(0)
}

func verify(url *string, username string, otp string) (bool, error) {
	requestStr := fmt.Sprintf(DEFAULT_REQUEST_FORMAT, username, otp)
	data := []byte(requestStr)
	if *IS_DEBUG {
		log.Println("POSTing:", requestStr)
	}

	req, err := http.NewRequest("POST", *url, bytes.NewBuffer(data))
	if err != nil {
		log.Println("Error reading request. ", err)
		return false, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Host", "localhost")

	// Create and Add cookie to request
	//cookie := http.Cookie{Name: "cookie_name", Value: "cookie_value"}
	//req.AddCookie(&cookie)

	// Set client opts
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var client = &http.Client{
		Timeout: 10 * time.Second,
	}

	if strings.HasPrefix(*DEFAULT_OTP_URL, "https://") {
		client.Transport = customTransport
	}

	// Validate cookie and headers are attached

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error reading response. ", err)
		return false, err
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)

	respVal := Resp{}
	if err := json.NewDecoder(resp.Body).Decode(&respVal); err != nil {
		log.Println(err)
		return false, err
	}
	return respVal.Valid, nil
}
