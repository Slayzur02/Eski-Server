package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Slayzur02/GoChess/pkg/constants"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CheckSignUpInfo(u UserCreation) (int, error)
	AddUser(u *UserCreation) (int, error)
	CheckLogInInfo(username string, hashPwd string) (int, error)
	InsertCookie(cookie *http.Cookie, userID int) error
	CheckCookieExists(cookieVal string) (int, error)
	UpdateCookieExp(cookieVal string) error
}

type Service struct {
	r Repository
}

func NewService(repo Repository) Service {
	return Service{
		r: repo,
	}
}

// hashing helper for passwords
func hashAndSalt(pwd string) string {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

// Adds user (no check)
func (s *Service) addUser(u UserCreation) error {
	// hash before adding
	u.Password = hashAndSalt((u.Password))
	id, err := s.r.AddUser(&u)
	if err != nil {
		return err
	}
	fmt.Println("CREATED USER WITH ID", id)
	return nil
}

// sign up process
func (s *Service) SignUp(u UserCreation) error {
	// make suure credentials are okay
	statusInt, err := s.r.CheckSignUpInfo(u)
	if err != nil {
		return err
	}

	if statusInt == 1 {
		return errors.New("email already exists")
	} else if statusInt == 2 {
		return errors.New("username already exists")
	}

	// add user
	err = s.addUser(u)
	if err != nil {
		return err
	}

	return nil
}

// logs in the user and returns new cookie/ error
func (s *Service) LogIn(l LogInCredentials) (*http.Cookie, error) {

	id, err := s.r.CheckLogInInfo(l.Username, hashAndSalt(l.Password))

	// bad credentials
	if err.Error() == "invalid log in credentials" {
		return nil, errors.New("invalid credentials")
	}

	// actual error
	if err != nil {
		return nil, err
	}

	// generate random cookie with UUID
	randomCookie, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	// create cookie
	cookie := http.Cookie{
		Name:     "auth",
		Value:    randomCookie.String(),
		Expires:  time.Now().Add(constants.CookieHours * time.Hour),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	}

	// insert cookie
	s.r.InsertCookie(&cookie, id)

	return &cookie, nil
}

// validates user request based on cookie - returns new cookie, userId /err
func (s *Service) AuthUserRequest(ck *http.Cookie) (*http.Cookie, int, error) {
	if time.Now().Sub(ck.Expires) > 0 {
		return nil, -1, errors.New("cookie has expired")
	}

	// get id from cookie and check db
	userID, err := s.r.CheckCookieExists(ck.Value)
	if err != nil {
		return nil, -1, err
	}

	// success, so let's update the expiration date
	ck.Expires = time.Now().Add(24 * 365 * time.Hour)

	// update cookie in-memory
	err = s.r.UpdateCookieExp(ck.Value)
	if err != nil {
		return nil, -1, err
	}

	return ck, userID, nil
}
