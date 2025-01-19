package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/DiscoDoggy/terabytes/go_backend/docs"
	"github.com/DiscoDoggy/terabytes/go_backend/internal/store"

)

//where PAI lives

//config and parameters injected into application
type application struct {
	config 	config
	store 	store.Storage
	logger *zap.SugaredLogger
}

type config struct {
	serverAddr 	string
	db			dbConfig
	env			string
	apiURL 		string
	mail		mailConfig
}

type mailConfig struct {
	exp time.Duration
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
			r.Put("/activate/{token_id}", app.activateUserHandler)
			r.Route("/{user_id}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)

				r.Get("/", app.getUserByIdHandler)
				r.Get("/feed", app.getUserFeedHandler)
				r.Put("/follow", app.followUserHandler)
				r.Delete("/unfollow", app.unfollowUserHandler)
			})

		})
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.CreateUserHandler)
		})
	}) 

	return r
}

func (app *application) run(mux http.Handler) error {
	//docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL 
	docs.SwaggerInfo.BasePath = "/v1"
	server := http.Server{
		Addr: ":8000",
		Handler: mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Minute,
	}

	app.logger.Infow("Server has started","addr", server.Addr, "env", app.config.env)

	return server.ListenAndServe()
}