package main

import (
	"log"

	"github.com/DiscoDoggy/terabytes/go_backend/internal/db"
	"github.com/DiscoDoggy/terabytes/go_backend/internal/store"
)

func main() {
	addr :="postgres://root:!Sti64fri@terabytes-db.c7gc2yii6ymp.us-west-1.rds.amazonaws.com:5432/postgres"

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store, conn)
}