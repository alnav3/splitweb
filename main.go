package main

import (
	"front"
	"log"
	"net/http"
	"github.com/a-h/templ"
)

func main() {
    // handle paths
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        c := front.LoginPage()
        templ.Handler(c).ServeHTTP(w, r)
	})

    // handle directories for static files
    handleDirectory( "/js/", "/style/" )

    // start the server
	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleDirectory(dirs... string) {
    for _, dir := range dirs {
        http.Handle(dir, http.StripPrefix(dir, http.FileServer(http.Dir("."+dir))))
    }
}
