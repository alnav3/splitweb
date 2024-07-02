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
        handleHome(w, r)
	})

    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        c := front.LoginBase(front.LoginBox(""))
        templ.Handler(c).ServeHTTP(w, r)
    })

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
        authUser(w, r)
	})

	http.HandleFunc("/signUp", func(w http.ResponseWriter, r *http.Request) {
        signupUser(w, r)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        c := front.LoginBase(front.RegisterBox(""))
        templ.Handler(c).ServeHTTP(w, r)
	})

	http.HandleFunc("/retry", func(w http.ResponseWriter, r *http.Request) {
        retry(w, r)
	})

    // handle directories for static files
    handleDirectory( "/js/", "/style/" , "/img/")
}

func retry(w http.ResponseWriter, r *http.Request) {
    sessions, err := store.Get(r, "cred")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    token := sessions.Values["token"]
    fmt.Println("token from retry:", token)
    c := front.LoginBase(front.LoginBox(""))
    templ.Handler(c).ServeHTTP(w, r)
}

func authUser(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")
    session, _ := store.Get(r, "cred")

    token, err := db.AuthUser(username, password, database)
    if err != nil {
        changeUrl(w, "/login")
        c := front.LoginBase(front.LoginBox("Invalid credentials"))
        templ.Handler(c).ServeHTTP(w, r)
        return
    }
    fmt.Println("token from authUser:", token)

    session.Values["token"] = token
    err = session.Save(r, w)

    if err != nil {
        changeUrl(w, "/login")
        c := front.LoginBase(front.LoginBox("Connection error. Try again later."))
        templ.Handler(c).ServeHTTP(w, r)
        return
    }

    changeUrl(w, "/")
    c := front.Base(username)
    templ.Handler(c).ServeHTTP(w, r)
}


func signupUser(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")
    passwordR := r.FormValue("passwordR")

    if password != passwordR {
        log.Println("Passwords do not match")
        c := front.LoginBase(front.RegisterBox("Passwords do not match"))
        templ.Handler(c).ServeHTTP(w, r)
        return
    }

    password, err := encryption.Encrypt(password)
    if err != nil {
        c := front.LoginBase(front.RegisterBox("Error encrypting password"))
        templ.Handler(c).ServeHTTP(w, r)
        return
    }

    err = db.SignupUser(username, password, database)
    if err != nil {
        c := front.LoginBase(front.RegisterBox("Error signing up"))
        templ.Handler(c).ServeHTTP(w, r)
        return
    }

    authUser(w, r)
}

func changeUrl(w http.ResponseWriter, url string) {
    w.Header().Set("Hx-Replace-Url", "true")
    w.Header().Set("Hx-Push-Url", url)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    sessions, err := store.Get(r, "cred")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    token, ok := sessions.Values["token"].(string)
    if !ok {
        // redirect to login
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    username, err := encryption.ValidateToken(token, db.GetSecret())
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    } else {
        c := front.Base(username)
        templ.Handler(c).ServeHTTP(w, r)
    }
}


func handleDirectory(dirs... string) {
    for _, dir := range dirs {
        http.Handle(dir, http.StripPrefix(dir, http.FileServer(http.Dir("."+dir))))
    }
}

func throwInternalServerError(w http.ResponseWriter, err error) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
}

