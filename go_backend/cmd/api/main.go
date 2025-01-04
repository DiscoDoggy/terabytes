package main

import(
	"log"
	"github.com/DiscoDoggy/terabytes/go_backend/internal/env"
)

func main() {

	var cfg = env.InitConfig()
	
	app := &application {
		config: cfg,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}