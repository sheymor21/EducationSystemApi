package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	port uint
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var conf config
	flag.UintVar(&conf.port, "port", 8080, "port to listen on")
	flag.StringVar(&conf.env, "env", "development", "environment to use dev|prod|test")
	flag.Parse()
	addr := fmt.Sprintf(":%d", conf.port)

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	app := &application{
		config: conf,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting server on %s in %s environmnent", addr, conf.env)
	err := srv.ListenAndServe()
	if err != nil {
		return
	}

}
