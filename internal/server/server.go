package server

import (
	"SchoolManagerApi/internal/database"
	"SchoolManagerApi/internal/utilities"
	"SchoolManagerApi/internal/validations"
	"flag"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"net/http"
	"os"
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

	envErr := godotenv.Load(".env")
	if envErr != nil {
		utilities.Log.Warnln(".env file not found")
	}

	if mc.Username == "" && mc.Password == "" {
		key := os.Getenv("SECRET_KEY")
		if key != "" {
			validations.SetSecretKey([]byte(key))
		} else {
			utilities.Log.Fatalln("Empty Secret Key")
		}

		mc.Username = os.Getenv("DB_U")
		mc.Password = os.Getenv("DB_P")
		dbName := os.Getenv("DB_NAME")
		dbUri := os.Getenv("DB_URI")
		if dbName != "" {
			mc.DbName = dbName
		}
		if dbUri != "" {
			mc.DbUri = dbUri
		}
	}

	database.SetMongoConfig(mc)
	database.Run()
	dbContext := database.GetMongoContext()
	defer database.CloseConnection(dbContext.Client)
	addr := fmt.Sprintf(":%d", conf.port)

	toURL, fileErr := utilities.FilePathToURL("./docs/swagger.json")
	if fileErr != nil {
		utilities.Log.Errorln(fileErr)
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

	utilities.Log.Infof("Starting server on %s in %s environment", addr, conf.env)
	err := srv.ListenAndServe()
	if err != nil {
		utilities.Log.Fatalf("Error starting server: %s", err)
		return
	}
}
