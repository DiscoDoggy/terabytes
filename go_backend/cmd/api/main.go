package main

import (
	"fmt"

	"github.com/DiscoDoggy/terabytes/go_backend/internal/db"
	"github.com/DiscoDoggy/terabytes/go_backend/internal/env"
	"github.com/DiscoDoggy/terabytes/go_backend/internal/store"
	
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const version = "0.0.1"
//	@title			Terabytes API
//	@description	API for Terabytes tech blog social network.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {

	cfg := config{
		serverAddr: env.GetString("HOST", ":8000"),
		db: dbConfig{
			addr: AssembleDBURL(),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 5),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 5),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8000"),
	}

	//logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()
	//Database
	db, err := db.New(
		cfg.db.addr, 
		cfg.db.maxOpenConns, 
		cfg.db.maxIdleConns, 
		cfg.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Info("Database connection pool established")
	
	store := store.NewStorage(db)

	app := &application {
		config: cfg,
		store: store,
		logger: logger,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}

func AssembleDBURL() string {
	godotenv.Load("secrets.env")

	dbUname := env.GetString("DB_USERNAME", "root")
	dbUrl := env.GetString("DB_URL", "URL")
	dbPort := env.GetString("DB_PORT", "5432")
	dbName := env.GetString("DB_NAME", "postgres")
	dbPassword := env.GetString("DB_PASSWORD", "admin")

	dbConnectionLink := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUname, dbPassword, dbUrl, dbPort, dbName)

	return dbConnectionLink	
}