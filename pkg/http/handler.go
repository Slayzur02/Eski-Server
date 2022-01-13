package http

import (
	"net/http"
	"time"

	"github.com/Slayzur02/GoChess/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func HandleCreation(a auth.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(5 * time.Second))

	r.Use(middleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", signUpHandler(a))
		r.Post("/login", logInHandler(a))
	})

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeChessWebsocket(w, r)
	})

	return r
}
