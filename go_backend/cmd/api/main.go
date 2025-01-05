package main

import (
	"fmt"
	"log"

	"github.com/DiscoDoggy/terabytes/go_backend/internal/db"
	"github.com/DiscoDoggy/terabytes/go_backend/internal/env"
	"github.com/DiscoDoggy/terabytes/go_backend/internal/store"
	"github.com/joho/godotenv"
)

func main() {

	cfg := config{
		serverAddr: env.GetString("HOST", ":8000"),
		db: dbConfig{
			addr: assembleDBURL(),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 5),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 5),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	db, err := db.New(
		cfg.db.addr, 
		cfg.db.maxOpenConns, 
		cfg.db.maxIdleConns, 
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}
	
	store := store.NewStorage(db)

	app := &application {
		config: cfg,
		store: store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}

func assembleDBURL() string {
	godotenv.Load("secrets.env")

	dbUname := env.GetString("DB_USERNAME", "root")
	dbUrl := env.GetString("DB_URL", "URL")
	dbPort := env.GetString("DB_PORT", "5432")
	dbName := env.GetString("DB_NAME", "postgres")
	dbPassword := env.GetString("DB_PASSWORD", "admin")

	dbConnectionLink := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUname, dbPassword, dbUrl, dbPort, dbName)

	return dbConnectionLink	
}