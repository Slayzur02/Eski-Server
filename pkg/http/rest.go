package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Slayzur02/GoChess/pkg/auth"
)

func signUpHandler(a auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")

		user := auth.UserCreation{Email: email, Username: username, Password: password}
		token, err := a.GenerateJWTToken(email)
		if err != nil {
			fmt.Println("failed to send email verification")
			http.Error(w, "couldn't send email verification", 404)
		}

		err = a.SendVerificationEmail(email, token)
		if err != nil {
			fmt.Println("email server didn't work")
		}

		err = a.SignUp(user)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "an error occured", 404)
		}
	}
}

func logInHandler(a auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		cookie, err := a.LogIn(auth.LogInCredentials{Email: email, Password: password})
		if err != nil {
			// no verification
			if err.Error() == "no email verification" {
				http.SetCookie(w, cookie)
				w.Header().Set("Content-Type", "application/json")
				verify := map[string]bool{
					"verify": false,
				}
				json.NewEncoder(w).Encode(verify)
			}
			// actual error
			fmt.Println(err)
			http.Error(w, err.Error(), 404)
		}

		http.SetCookie(w, cookie)
		w.Header().Set("Content-Type", "application/json")
		verify := map[string]bool{
			"verify": true,
		}
		json.NewEncoder(w).Encode(verify)
	}
}

func logOutHandler(a auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ck, err := r.Cookie("auth")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), 404)
		}

		err = a.LogOut(ck)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), 404)
		}
	}
}

func logOutAllHandler(a auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ck, err := r.Cookie("auth")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), 404)
		}

		err = a.LogOutAll(ck)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), 404)
		}
	}
}

// log out specific instance: get cookie, remove in redis

// log out all instances: get cookie, get ID, remove all isntances with ID

// verification: sign up -> extra verify state in email, send email with random id, store in redis; add specific route on frontend ->
// verfication route: gets the randomly spawned id for the verification, check with redis, verify

// change password: send email with password reset route with random ID (maps to userID in redis); add specific route on frontend to fill
// in the new password. This route sends to /password/reset on backend along with ID, which changes if the ID exists in Redis.

// log in middleware
