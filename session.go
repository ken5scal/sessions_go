package main

import (
	"net/http"
	"time"
	"github.com/satori/go.uuid"
)

const CookieName = "session"
const sessionLength int = 30 // sec

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

func createNewSession() *http.Cookie {
	sessionId := uuid.NewV4()
	cookie := &http.Cookie{
		Name: CookieName,
		Value: sessionId.String(),
		MaxAge: sessionLength,
	}
	return cookie
}

func clearSessions() {
	for sessionId, session := range dbSessions {
		if time.Now().Sub(session.lastActivity) > (time.Duration(sessionLength) * time.Second) {
			delete(dbSessions, sessionId)
		}
	}
}

func getUser(res http.ResponseWriter, req *http.Request) user {
	cookie, err := req.Cookie(CookieName)
	if err != nil {
		cookie = createNewSession()
	}
	cookie.MaxAge = sessionLength
	http.SetCookie(res, cookie)

	var u user
	if session, ok := dbSessions[cookie.Value]; ok {
		session.lastActivity = time.Now()
		dbSessions[cookie.Value] = session
		u = dbUsers[session.username]
	}
	return u
}