package server

import (
	"calificationApi/internal/database"
	"calificationApi/internal/utilities"
	"flag"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type config struct {
	port uint
	env  string
}

type application struct {
	config         config
	swaggerSpecURL string
	validator      *validator.Validate
}

func ListenServer() {
	var conf config
	var mc database.MongoConfig
	flag.StringVar(&mc.DbName, "DB_NAME", "EducationSystem", "MongoDB NAME")
	flag.StringVar(&mc.DbUri, "DB_URI", "mongodb://localhost:27017", "MongoDB URI")
	flag.StringVar(&mc.Username, "DB_U", "", "MongoDB Username")
	flag.StringVar(&mc.Password, "DB_P", "", "MongoDB Password")
	flag.UintVar(&conf.port, "port", 8080, "port to listen on")
	flag.StringVar(&conf.env, "env", "dev", "environment to use dev|prod|test")
	flag.Parse()
	database.SetMongoConfig(mc)
	database.Run()
	dbContext := database.GetMongoContext()
	defer database.CloseConnection(dbContext.Client)
	addr := fmt.Sprintf(":%d", conf.port)

	toURL, fileErr := utilities.FilePathToURL("./docs/swagger.json")
	if fileErr != nil {
		utilities.Log.Fatal(fileErr)
		return
	}

	app := &application{
		config:         conf,
		swaggerSpecURL: toURL,
		validator:      validator.New(validator.WithRequiredStructEnabled()),
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	utilities.Log.Infof("Starting server on %s in %s environmnent", addr, conf.env)
	err := srv.ListenAndServe()
	if err != nil {
		return
	}
}
