package main

import (
	"net/http"
	"time"
)

const CookieName = "session"

func alreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
	c, err := req.Cookie(CookieName)
	if err != nil {
		return false
	}
	s, ok := dbSessions[c.Value];
	if ok {
		s.lastActivity = time.Now()
		dbSessions[c.Value] = s
	}

	_, ok = dbUsers[s.username]
	c.MaxAge = sessionLength
	http.SetCookie(w, c)
	return ok
}