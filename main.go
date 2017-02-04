package main

import (
	"net/http"
	"time"
	"html/template"
	"log"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
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

const port = ":8080"
const sessionLength int = 30 // sec
var tpl *template.Template
var dbUsers = make(map[string]user) //key -> user iD
var dbSessions = make(map[string]session)   //key -> sessionId(cookie value)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
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
	// Show User if request parameter contains session cookie

	tpl.ExecuteTemplate(w, "index.html", nil)
}

func bar(w http.ResponseWriter, req *http.Request) {

}

func signup(w http.ResponseWriter, req *http.Request) {
	// Check If user already signs in
	// if signs in, redirect to root dir
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
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
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name: CookieName,
			Value: sID.String(),
			MaxAge: sessionLength,
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = session{un, time.Now()}

		// Now redirect to Top pageu
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(w, "signup.html", u)
}

func login(w http.ResponseWriter, req *http.Request) {

}

func logout(w http.ResponseWriter, req *http.Request) {

}

