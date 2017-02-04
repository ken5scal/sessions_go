package main

import (
	"net/http"
	"time"
	"html/template"
	"log"
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
var dbSessions = make(map[string]session)   //key -> sessionId

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
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func bar(w http.ResponseWriter, req *http.Request) {

}

func signup(w http.ResponseWriter, req *http.Request) {

}

func login(w http.ResponseWriter, req *http.Request) {

}

func logout(w http.ResponseWriter, req *http.Request) {

}

