package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Slayzur02/GoChess/pkg/auth"
	"github.com/Slayzur02/GoChess/pkg/constants"
	sqlc "github.com/Slayzur02/GoChess/pkg/internal_sql"
	"golang.org/x/crypto/bcrypt"
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

func comparePlainAndHash(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))

	return (err == nil)
}

func (r *Repo) CheckLogInInfo(email string, plainPwd string) (int, error) {
	idPwdVerified, err := r.db.GetPwdIdVerifyfromEmail(context.Background(), email)
	if err != nil {
		log.Print(err)
		return -1, nil
	}

	fmt.Println(idPwdVerified.Hashedpwd, plainPwd, comparePlainAndHash(idPwdVerified.Hashedpwd, plainPwd))
	if !comparePlainAndHash(idPwdVerified.Hashedpwd, plainPwd) {
		return -1, errors.New("invalid log in credentials")
	}

	if !idPwdVerified.Verified {
		return int(idPwdVerified.ID), errors.New("no email verification")
	}

	return int(idPwdVerified.ID), nil
}

// HELPERS for redis parsing
func parseSessionID(s string) string {
	return "session:" + s
}

func parseSessionInsertTime(s string) string {
	return "time:" + s
}

func parseIDBlacklist(s string) string {
	return "blacklist:" + s
}

func (r *Repo) InsertCookie(cookie *http.Cookie, userID int) error {
	// set the ckID - usID
	err := r.mc.Set(context.Background(), parseSessionID(cookie.Value), strconv.Itoa(userID),
		time.Until(cookie.Expires)).Err()
	if err != nil {
		return err
	}

	// sets the ckID - insert time
	err = r.mc.Set(context.Background(), parseSessionInsertTime(cookie.Value), time.Now(),
		time.Until(cookie.Expires)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) DeleteCookie(cookie *http.Cookie) error {
	// delete the ckID - usID
	err := r.mc.Del(context.Background(), parseSessionID(cookie.Value))
	if err != nil {
		return err.Err()
	}

	// delete the ckID - insert time
	err = r.mc.Del(context.Background(), parseSessionInsertTime(cookie.Value))
	if err != nil {
		return err.Err()
	}

	return nil
}

func (r *Repo) BlackistID(cookie *http.Cookie) error {
	s := r.mc.Get(context.Background(), parseSessionID(cookie.Value)).String()

	err := r.mc.Set(context.Background(), parseIDBlacklist(s), time.Now(), 0).Err()
	return err
}

func (r *Repo) ValidSession(cookie *http.Cookie) (bool, error) {
	id := r.mc.Get(context.Background(), parseSessionID(cookie.Value)).String()
	blacklistTime, err := r.mc.Get(context.Background(), parseIDBlacklist(id)).Time()
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	if (blacklistTime != time.Time{}) { // not empty
		t, err := r.mc.Get(context.Background(), parseSessionInsertTime(cookie.Value)).Time()
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		if t.Before(blacklistTime) {
			return false, nil
		}
	}

	return true, nil

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

func (r *Repo) VerifyEmail(email string) error {
	err := r.VerifyEmail(email)
	if err != nil {
		return err
	}

	return nil
}
