package http

import (
	"fmt"
	"net/http"

	"github.com/Slayzur02/GoChess/pkg/auth"
)

func postmanTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		randoCookie := http.Cookie{Name: "randoCookie", Value: "what?"}
		http.SetCookie(w, &randoCookie)
	}
}

func signUpHandler(a *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")

		user := auth.UserCreation{Email: email, Username: username, Password: password}

		err := a.SignUp(user)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "an error occured", 404)
		}
	}
}

func logInHandler(a *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		cookie, err := a.LogIn(auth.LogInCredentials{Username: username, Password: password})
		if err != nil {
			http.Error(w, "an error occured", 404)
		}

		http.SetCookie(w, cookie)
	}
}
