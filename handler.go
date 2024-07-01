package main

import (
	"fmt"
	"front"
	"log"
	"net/http"

	"git.alnav.dev/alnav3/splitweb/api/db"
	"git.alnav.dev/alnav3/splitweb/api/encryption"
	"github.com/a-h/templ"
)

func handlePaths() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        c := front.LoginPage("")
        templ.Handler(c).ServeHTTP(w, r)
	})

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
        authUser(w, r)
	})

	http.HandleFunc("/signUp", func(w http.ResponseWriter, r *http.Request) {
        signupUser(w, r)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        c := front.RegisterPage("")
        templ.Handler(c).ServeHTTP(w, r)
	})

	http.HandleFunc("/retry", func(w http.ResponseWriter, r *http.Request) {
        retry(w, r)
	})

    // handle directories for static files
    handleDirectory( "/js/", "/style/" )
}

func retry(w http.ResponseWriter, r *http.Request) {
    sessions, err := store.Get(r, "cred")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    token := sessions.Values["token"]
    fmt.Println("token from retry:", token)
    c := front.LoginPage("")
    templ.Handler(c).ServeHTTP(w, r)
}

func authUser(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")
    session, _ := store.Get(r, "cred")

    token, err := db.AuthUser(username, password, database)
    if err != nil {
        c := front.LoginPage("Invalid credentials")
        templ.Handler(c).ServeHTTP(w, r)
        return
    }
    fmt.Println("token from authUser:", token)

    session.Values["token"] = token
    err = session.Save(r, w)

    if err != nil {
        c := front.LoginPage("Connection error. Try again later.")
        templ.Handler(c).ServeHTTP(w, r)
        return
    }

    c := front.Token(token)
    templ.Handler(c).ServeHTTP(w, r)
}

func throwInternalServerError(w http.ResponseWriter, err error) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
}

func signupUser(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")
    passwordR := r.FormValue("passwordR")

    if password != passwordR {
        log.Println("Passwords do not match")
        c := front.RegisterPage("Passwords do not match")
        templ.Handler(c).ServeHTTP(w, r)
        return
    }

    password, err := encryption.Encrypt(password)
    if err != nil {
        c := front.RegisterPage("Error encrypting password")
        templ.Handler(c).ServeHTTP(w, r)
        return
    }

    err = db.SignupUser(username, password, database)
    if err != nil {
        c := front.RegisterPage("The user already exists")
        templ.Handler(c).ServeHTTP(w, r)
        return
    }
    authUser(w, r)
}


func handleDirectory(dirs... string) {
    for _, dir := range dirs {
        http.Handle(dir, http.StripPrefix(dir, http.FileServer(http.Dir("."+dir))))
    }
}
