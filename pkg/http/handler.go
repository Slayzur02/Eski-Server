package http

import (
	"net/http"

	"github.com/Slayzur02/GoChess/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func HandleCreation(a *auth.Service) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		r.Get("/postmanTest", postmanTest())
		r.Post("/signup", signUpHandler(a))
		r.Post("/login", logInHandler(a))
	})

	return r
}
