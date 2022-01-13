package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/Slayzur02/GoChess/pkg/constants"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CheckSignUpInfo(u UserCreation) (int, error)
	AddUser(u *UserCreation) (int, error)
	CheckLogInInfo(email string, hashPwd string) (int, error)
	InsertCookie(cookie *http.Cookie, userID int) error
	DeleteCookie(cookie *http.Cookie) error
	CheckCookieExists(cookieVal string) (int, error)
	UpdateCookieExp(cookieVal string) error
	VerifyEmail(email string) error
	BlackistID(ck *http.Cookie) error
}

type Service interface {
	SignUp(u UserCreation) error
	LogIn(l LogInCredentials) (*http.Cookie, error)
	AuthUserRequest(ck *http.Cookie) (*http.Cookie, int, error)
	GenerateJWTToken(email string) (string, error)
	SendVerificationEmail(email string, url string) error
	LogOut(ck *http.Cookie) error
	LogOutAll(ck *http.Cookie) error
}

type service struct {
	r Repository
}

func NewService(repo Repository) Service {
	return &service{
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
func (s *service) addUser(u UserCreation) error {
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
func (s *service) SignUp(u UserCreation) error {
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
func (s *service) LogIn(l LogInCredentials) (*http.Cookie, error) {

	id, err := s.r.CheckLogInInfo(l.Email, l.Password)

	if err != nil {
		// bad credentials
		if err.Error() == "invalid log in credentials" {
			return nil, errors.New("invalid credentials")
		}

		if err.Error() == "no email verification" {
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

			return &cookie, errors.New("no email verification")

		}

		// actual error
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
func (s *service) AuthUserRequest(ck *http.Cookie) (*http.Cookie, int, error) {
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

func (s *service) GenerateJWTToken(email string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", err
	}

	signKey := []byte(os.Getenv("secretJwtKey"))

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * constants.VerificationMinutes).Unix()

	tokenStr, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *service) VerifyWithJwt(tokenString string) error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	signKey := []byte(os.Getenv("secretJwtKey"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return signKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		s.r.VerifyEmail(claims["email"].(string))
		return nil
	} else {
		return err
	}
}

func (s *service) SendVerificationEmail(email string, url string) error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	from := os.Getenv("hostEmail")
	password := os.Getenv("hostPassword")
	smtpHost := "smtp.gmail.com"

	to := []string{
		"adt008@bucknell.edu",
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)
	message := []byte("I like you pepeHands message\r\n" + url + "\r\n")

	// trust that this works
	go smtp.SendMail(smtpHost+":587", auth, from, to, message)

	return nil
}

func (s *service) LogOut(ck *http.Cookie) error {
	err := s.r.DeleteCookie(ck)
	return err
}

func (s *service) LogOutAll(ck *http.Cookie) error {
	err := s.r.BlackistID(ck)
	return err
}
