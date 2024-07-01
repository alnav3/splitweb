package main

import (
	"database/sql"
	"log"
	"net/http"

    "github.com/gorilla/sessions"
	"git.alnav.dev/alnav3/splitweb/api/db"
)

var store = sessions.NewCookieStore([]byte(db.GetSessionToken()))
var database *sql.DB

func main() {
    // connect to the database
    var err error
    database, err = db.Connect()

    if err != nil {
        log.Fatal("Error connecting to database: ", err)
    }
    handlePaths()

    // start the server
	log.Println("Server started on port 80")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
        defer database.Close()
		log.Fatal("ListenAndServe: ", err)
	}
}

