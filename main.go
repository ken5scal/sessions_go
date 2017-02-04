package main

import (
	"net/http"
	"time"
	"html/template"
	"log"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

type user struct {
	UserName string
	Password []byte
	First string
	Last string
	Role string
}

type session struct {
	username string
	lastActivity time.Time
}

const port = ":1080"

var tpl *template.Template
var dbUsers = make(map[string]user) //key -> user iD
var dbSessions = make(map[string]session)   //key -> sessionId(cookie value)
var dbInitiatedTime time.Time

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	dbInitiatedTime = time.Now()
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatalln(http.ListenAndServe(port, nil))
}

func index(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	tpl.ExecuteTemplate(w, "index.html", u)
}

func bar(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "bar.html", u)
}

func signup(w http.ResponseWriter, req *http.Request) {
	// Check If user already signs in
	// if signs in, redirect to root dir
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	var u user
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")
		r := req.FormValue("role")

		if _, ok := dbUsers[un]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}

		// Then create and store user in User DB
		passwordInByte, err := bcrypt.GenerateFromPassword([]byte(p),bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed creating user", http.StatusInternalServerError)
			return
		}
		u = user{un, passwordInByte, f, l, r}
		dbUsers[un] = u

		// Create session and store in Session Db
		c := createNewSession()
		http.SetCookie(w,c)
		dbSessions[c.Value] = session{un, time.Now()}

		// Now redirect to Top pageu
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(w, "signup.html", u)
}

func login(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	var u user
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")

		u, ok := dbUsers[un]
		if !ok {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		cookie := createNewSession()
		http.SetCookie(w, cookie)

		dbSessions[cookie.Value] = session{un, time.Now()}
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "login.html", u)
}

func logout(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		fmt.Println("you are not logged in")
		return
	}

	cookie, _ := req.Cookie(CookieName)
	// delete the session
	delete(dbSessions, cookie.Value)
	// remove the cookie
	cookie = &http.Cookie{
		Name:   CookieName,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, nil)
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func ShowAllUsers() {
	for _, value :=  range dbSessions {
		fmt.Println()
		fmt.Println(value)
	}

	for _, value := range dbUsers {
		fmt.Println(value.UserName)
	}
}
