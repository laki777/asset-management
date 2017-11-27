package main

import (
	"asset-management/boltfunc"
	"net/http"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

type Configuration struct {
	sessionLength int
}

const sessionLength int = 3600

func setSession(user string, w http.ResponseWriter, req *http.Request) {
	cookie, errors := req.Cookie("session")
	if errors != nil {
		id, _ := uuid.NewV4()
		cookie = &http.Cookie{
			Name:  "session",
			Value: id.String(),
		}
		cookie.MaxAge = sessionLength
		http.SetCookie(w, cookie)
	}

	input := []*boltfunc.Session{
		{user, time.Now()},
	}

	for _, p := range input {
		p.SessionSave(cookie.Value)
	}
}

func alreadyLoggedIn(req *http.Request) bool {
	if !boltfunc.OpenDB {
		return false
	}
	cookie, errors := req.Cookie("session")
	if errors != nil {
		return false
	}
	p, error := boltfunc.SessionGetPerson(cookie.Value)
	if error != nil {
		return false
	}

	// boltfunc.List("sessions")
	// boltfunc.List("adminpeople")

	_, error = boltfunc.AdminGetPerson(p.Username, "adminpeople")
	if error != nil {
		return false
	}

	// refresh session ??? or expired ???

	return true
}

func cleanSessions(w http.ResponseWriter, req *http.Request) bool {
	cookie, errors := req.Cookie("session")
	if errors != nil {
		return false
	}

	error := boltfunc.SessionDelete(cookie.Value)
	if error != nil {
		return false
	}

	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	return true

}
