package main

import (
	"fmt"
	// "go/doc"
	"log"
	"net/http"
	"time"

	"github.com/DiscoDoggy/terabytes/go_backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggo/http-swagger"
	"github.com/DiscoDoggy/terabytes/go_backend/docs"
)

//where PAI lives

//config and parameters injected into application
type application struct {
	config 	config
	store 	store.Storage
}

type config struct {
	serverAddr 	string
	db			dbConfig
	env			string
	apiURL 	string

}

type dbConfig struct {
	addr 			string
	maxOpenConns	int
	maxIdleConns 	int
	maxIdleTime 	string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.serverAddr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/blogs", func(r chi.Router) {
			r.Post("/", app.createBlogHandler)
			r.Route("/{blog_id}", func(r chi.Router) {
				r.Use(app.blogPostContextMiddleware)
				
				r.Get("/", app.getBlogByIdHandler)
				r.Delete("/", app.deleteBlogByIdHandler)
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Route("/{user_id}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)

				r.Get("/", app.getUserByIdHandler)
				r.Get("/feed", app.getUserFeedHandler)
				r.Put("/follow", app.followUserHandler)
				r.Delete("/unfollow", app.unfollowUserHandler)
			})

		})
	}) 

	return r
}

func (app *application) run(mux http.Handler) error {
	//docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL 
	docs.SwaggerInfo.BasePath = "v1"
	server := http.Server{
		Addr: ":8000",
		Handler: mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Minute,
	}

	log.Printf("Server has started at %s", server.Addr)

	return server.ListenAndServe()
}