package router

import (
	"net/http"
	"video-streaming/handlers"

	"github.com/go-chi/chi/v5"
)

func RouteHandler() http.Handler {

	mux := chi.NewRouter()

	mux.Post("/user", handlers.InsertUser)
	mux.Post("/login", handlers.LoginUser)

	return mux
}
