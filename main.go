package main

import "net/http"

var port = ":8080"

func init() {

}

func main() {
	http.HandlerFunc("/", index)
	http.HandlerFunc("/bar", bar)
	http.HandlerFunc("/signup", signup)
	http.HandlerFunc("/login", login)
	http.HandlerFunc("/logout", logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe("port", nil)
}

func index(w http.Response, req *http.Request) {

}

func bar(w http.Response, req *http.Request) {

}

func signup(w http.Response, req *http.Request) {

}

func login(w http.Response, req *http.Request) {

}

func logout(w http.Response, req *http.Request) {

}

