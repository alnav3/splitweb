package main

import (
	"database/sql"
	"fmt"
	"front"
	"log"
	"net/http"

	"git.alnav.dev/alnav3/splitweb/api/db"
	"git.alnav.dev/alnav3/splitweb/api/encryption"
	"github.com/a-h/templ"
)

var database *sql.DB

func main() {
    // connect to the database
    var err error
    database, err = db.Connect()

    handleErr(err, "Error connecting to database: ")
    handlePaths()

    // start the server
	log.Println("Server started on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
        defer database.Close()
		log.Fatal("ListenAndServe: ", err)
	}
}

func handlePaths() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        c := front.LoginPage()
        templ.Handler(c).ServeHTTP(w, r)
	})

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
        authUser(w, r)
	})

    // handle directories for static files
    handleDirectory( "/js/", "/style/" )
}

func authUser(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")

    token, err := db.AuthUser(username, password, database)
    if err != nil {
        panic(fmt.Sprintf("Error authenticating user: ", err))
    }
    fmt.Println(token)

    c := front.Token(token)
    templ.Handler(c).ServeHTTP(w, r)
}
func signupUser(w http.ResponseWriter, r *http.Request) {
    fmt.Println("holi")
    username := r.FormValue("username")
    password := r.FormValue("password")
    password, err := encryption.Encrypt(password)
    if err != nil {
        log.Println("Error encrypting password: ", err)
        return
    }
    err = db.SignupUser(username, password, database)
    c := front.LoginPage()
    templ.Handler(c).ServeHTTP(w, r)
}


func handleDirectory(dirs... string) {
    for _, dir := range dirs {
        http.Handle(dir, http.StripPrefix(dir, http.FileServer(http.Dir("."+dir))))
    }
}

func handleErr(err error, msg string) {
    if err != nil {
        log.Fatal(msg, err)
    }
}
