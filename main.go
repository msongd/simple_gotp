package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	CONFIG_FILE = flag.String("f", "config.json", "Config file")
	GLOBAL_CFG  = NewConfig()
)

type Env struct {
	Db *OtpConfig
}

func globalInit() {
	if len(os.Args) == 1 {
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := GLOBAL_CFG.LoadConfig(*CONFIG_FILE)
	if err != nil {
		log.Fatal("Load config err:", err)
	}
	fmt.Printf("%v\n", *GLOBAL_CFG)
	log.SetOutput(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/app.log", GLOBAL_CFG.LogDir),
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	})
}
func main() {
	globalInit()
	log.Println("Starting")
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	workingEnv := &Env{}
	db, err := LoadFromFile(GLOBAL_CFG.DataFile)
	if err != nil {
		log.Fatal("Unable to load data file, err:", err)
	}
	workingEnv.Db = db
	defer workingEnv.Db.SaveToFile(GLOBAL_CFG.DataFile)
	r := MakeRouter(workingEnv)
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
		Handler:      handlers.CombinedLoggingHandler(l, r), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("Shutting down")
	os.Exit(0)
}
func test() {
	cfg := NewOtpConfig()
	cfg.AddToken("tungdm@pam", "pam")
	cfg.AddToken("tungdm@pam", "pam")
	cfg.AddToken("khanhtn@pam", "pam")
	cfg.AddToken("duyhd@pam", "pam")
	cfg.AddToken("khanhtn@pam", "pam")
	cfg.SetActiveToken("khanhtn@pam", "2")
	cfg.SetActiveToken("duyhd@pam", "3")
	cfg.SetActiveToken("xxx", "0")
	cfg.Dump()
	found, a1 := cfg.GetActiveToken("khanhtn@pam")
	fmt.Println(found, a1)
	found, a2 := cfg.GetActiveToken("duyhd@pam")
	fmt.Println(found, a2)
	found, a3 := cfg.GetActiveToken("xx")
	fmt.Println(found, a3)
	cfg.RemoveToken("khanhtn@pam", "1")
	cfg.AddToken("khanhtn@pam", "pam")
	cfg.Dump()
}
