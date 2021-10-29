package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	CONFIG_FILE   = flag.String("f", "config.json", "Config file")
	DEBUG_FLAG    = flag.Bool("d", false, "Turn on more debug messages")
	DUMP_CFG_ONLY = flag.Bool("t", false, "Dump config and exit")
	GLOBAL_CFG    = NewConfig()
	WAIT          time.Duration
)

type Env struct {
	Db  *OtpConfig
	Cfg *Config
}

func globalInit() {
	if len(os.Args) == 1 {
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.DurationVar(&WAIT, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := GLOBAL_CFG.LoadConfig(*CONFIG_FILE)
	if err != nil {
		log.Fatal("Load config err:", err)
	}
	//fmt.Printf("%v\n", *GLOBAL_CFG)
	log.SetOutput(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/app.log", GLOBAL_CFG.LogDir),
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	})
	if GLOBAL_CFG.KeycloakCfg == nil || GLOBAL_CFG.KeycloakCfg.JwkUrl == "" {
		log.Println("Empty JWT url")
		GLOBAL_CFG.NoAuth = true
	}
	if GLOBAL_CFG.SSLKeyFile != "" {
		/*
			// for client cert verify (mutual SSL)
			err = GLOBAL_CFG.LoadTLS()
			if err != nil {
				log.Fatal("Unable to load certificates")
			}
		*/
		_, err = ioutil.ReadFile(GLOBAL_CFG.SSLCertFile)
		if err != nil {
			log.Fatal("Unable to read cert file:", err)
		}
		_, err = ioutil.ReadFile(GLOBAL_CFG.SSLKeyFile)
		if err != nil {
			log.Fatal("Unable to read key file:", err)
		}
	}
	if *DUMP_CFG_ONLY {
		fmt.Printf("%+v\n", GLOBAL_CFG)
		os.Exit(0)
	}
}
func main() {
	globalInit()
	if *DEBUG_FLAG {
		log.Println("Starting")
		log.Printf("%+v\n", GLOBAL_CFG)
	}

	//var wait time.Duration
	//
	//flag.Parse()

	workingEnv := &Env{}
	db, err := LoadFromFile(GLOBAL_CFG.DataFile)
	if err != nil {
		log.Fatal("Unable to load data file, err:", err)
	}
	if *DEBUG_FLAG {
		log.Printf("%+v\n", db)
	}
	workingEnv.Db = db
	workingEnv.Cfg = GLOBAL_CFG
	defer workingEnv.Db.SaveToFile(GLOBAL_CFG.DataFile)
	r := MakeRouter(workingEnv)
	headersOk := handlers.AllowedHeaders([]string{"*"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	l := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/access.log", GLOBAL_CFG.LogDir),
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	}
	srv := &http.Server{
		Addr: GLOBAL_CFG.Listen,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handlers.CORS(headersOk, originsOk, methodsOk)(handlers.CombinedLoggingHandler(l, r)), // Pass our instance of gorilla/mux in.
	}
	/* for client ssl verification (mutual ssl)
	if workingEnv.Cfg.SSLKeyFile != "" {
		srv.TLSConfig = workingEnv.Cfg.TLSConfig
	}*/
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if workingEnv.Cfg.SSLKeyFile != "" {
			if err := srv.ListenAndServeTLS(workingEnv.Cfg.SSLCertFile, workingEnv.Cfg.SSLKeyFile); err != nil {
				log.Println("Start SSL server err:", err)
			}
		} else if err := srv.ListenAndServe(); err != nil {
			log.Println("Start server err:", err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), WAIT)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	//workingEnv.Db.SaveToFile(GLOBAL_CFG.DataFile)
	log.Println("Shutting down")
	//os.Exit(0)
}
