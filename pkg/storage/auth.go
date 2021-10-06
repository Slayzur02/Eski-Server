package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Slayzur02/GoChess/pkg/auth"
	"github.com/Slayzur02/GoChess/pkg/constants"
	sqlc "github.com/Slayzur02/GoChess/pkg/internal_sql"
)

func (r *Repo) AddUser(a *auth.UserCreation) (int, error) {
	newId, err := r.db.InsertUser(context.Background(), sqlc.InsertUserParams{Email: a.Email, Username: a.Username, Hashedpwd: a.Password})
	if err != nil {
		return -1, err
	}
	return int(newId), nil
}

func (r *Repo) CheckSignUpInfo(u auth.UserCreation) (int, error) {
	_, err := r.db.GetIDFromEmail(context.Background(), u.Email)

	if err == sql.ErrNoRows {
		_, err = r.db.GetIDfromUsername(context.Background(), u.Username)
		if err == sql.ErrNoRows {
			return 0, nil // success!
		}
		if err != nil {
			return -1, err // some error other than no rows occured
		}

		// err is nil, username exists
		return 2, err
	}

	if err != nil {
		return -1, err // some error other than no rows occured
	}

	// err is nil, email exists
	return 1, nil
}

func (r *Repo) CheckLogInInfo(username string, hashedPwd string) (int, error) {
	idAndPwd, err := r.db.GetHashedPwdAndIdfromUsername(context.Background(), username)
	if err != nil {
		log.Print(err)
		return -1, nil
	}

	if idAndPwd.Hashedpwd != username {
		return -1, errors.New("invalid log in credentials")
	}

	return int(idAndPwd.ID), nil
}

func (r *Repo) InsertCookie(cookie *http.Cookie, userID int) error {
	err := r.mc.Set(context.Background(), cookie.Value, strconv.Itoa(userID),
		time.Until(cookie.Expires)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) CheckCookieExists(cookieVal string) (int, error) {
	id, err := r.mc.Get(context.Background(), cookieVal).Result()
	if err != nil {
		return -1, nil
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return -1, err
	}
	return idInt, nil
}

func (r *Repo) UpdateCookieExp(cookieVal string) error {
	duration := time.Until(time.Now().Add(time.Hour * constants.CookieHours))
	_, err := r.mc.Expire(context.Background(), cookieVal, duration).Result()
	if err != nil {
		return err
	}
	return nil

}
