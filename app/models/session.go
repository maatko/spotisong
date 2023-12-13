package models

import (
	"errors"
	"net/http"
	"os"
	"spotisong/api"
	"strconv"
	"time"
)

type Session struct {
	ID        int  `key:"primary"`
	User      User `key:"foreign"`
	CreatedAt time.Time
	ExpiresAt time.Time
}

const AppSessionCookie = "session_id"

////////////////////////////////////////
// Utility functions
////////////////////////////////////////

func NewSession(user User) *Session {
	expiresIn, err := strconv.Atoi(os.Getenv("APP_SESSION_DURATION"))
	if err != nil {
		panic(err)
	}

	currentTime := api.TimeCurrent()
	return &Session{
		User:      user,
		CreatedAt: currentTime,
		ExpiresAt: currentTime.Add(time.Second * time.Duration(expiresIn)),
	}
}

func GetCookieSession(request *http.Request) (*Session, error) {
	cookie, err := api.AppCookieStore.Get(request, AppSessionCookie)
	if err != nil {
		return nil, err
	}

	value, ok := cookie.Values[AppSessionCookie]
	if !ok {
		return nil, errors.New("invalid session")
	}

	session, ok := value.(Session)
	if !ok {
		return nil, errors.New("invalid session")
	}

	return &session, nil
}

func NewCookieSession(request *http.Request, response http.ResponseWriter, user User) (*Session, error) {
	cookie, err := api.AppCookieStore.Get(request, AppSessionCookie)
	if err != nil {
		return nil, err
	}

	session := NewSession(user)

	err = session.Save()
	if err != nil {
		session = nil
		return nil, err
	}

	cookie.Values[AppSessionCookie] = session

	err = cookie.Save(request, response)
	if err != nil {
		session = nil
		return nil, err
	}

	return session, nil
}

////////////////////////////////////////
// Database managing functions
////////////////////////////////////////

func (session *Session) Load(keys ...string) error {
	return api.FetchModel(session, keys...)
}

func (session *Session) Save() error {
	id, err := api.SaveModel(*session)
	if err != nil {
		return err
	}

	session.ID = int(id)
	return nil
}
